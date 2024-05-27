/***********************************************************************
TV Controller
Copyright 2024 by Volodymyr Dobryvechir (vdobryvechir@gmail.com)
************************************************************************/

package tvcontrol

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Dobryvechir/microcore/pkg/dvdbmanager"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
)

type TaskWorker struct {
	Id            string
	Task          *TvTask
	WakeUpChannel chan int
	StopChannel   chan int
}

func (task *TaskWorker) RunBackground() {
	err := task.LoadTask()
	if err != nil {
		dvlog.PrintError(err)
		return
	}
	delay := GetDelayInIdleCase()
	for {
		res, err := task.RunNextTask()
		if err != nil {
			dvlog.PrintError(err)
			delay = GetDelayInErrorCase()
		} else if res {
			delay = GetDelayInOperationCase()
		} else {
			delay = GetDelayInIdleCase()
		}
		timer := time.NewTimer(time.Duration(delay) * time.Second)
		select {
		case val := <-task.StopChannel:
			if logLevel {
				dvlog.PrintfFullOnly("b worker %s stopped %d", task.Id, val)
			}
			close(task.StopChannel)
			close(task.WakeUpChannel)
			return
		case wval := <-task.WakeUpChannel:
			err = task.LoadTask()
			if logLevel || err != nil {
				dvlog.PrintfFullOnly("b worker %s waken up %d %v", task.Id, wval)
			}
		case <-timer.C:
			if logLevel {
				dvlog.PrintfFullOnly("Continue to work by timer %d", delay)
			}
		}
	}
}

func (task *TaskWorker) LoadTask() error {
	res, err := dvdbmanager.RecordReadOne(taskDbName, task.Id)
	if err != nil {
		return err
	}
	if res == nil {
		return errors.New("task no longer exists")
	}
	task.Task = &TvTask{}
	err = res.DvVariableToAnyStruct(task.Task)
	return err
}

func (task *TaskWorker) RunNextTask() (bool, error) {
	if task == nil || task.Task == nil || len(task.Id) == 0 || len(task.Task.Url) == 0 {
		return false, nil
	}
	t := task.Task
	if len(t.NewPresentationId) == 0 || len(t.NewPresentationVersion) == 0 {
		return false, task.RunCheckConnection()
	}
	if t.NewPresentationId != t.OldPresentationId || t.NewPresentationVersion != t.OldPresentationVersion {
		return true, task.RunConfigSending()
	}
	if len(t.LeftFiles) != 0 {
		return true, task.RunFileSending()
	}
	return false, task.RunCheckConnection()
}

func (task *TaskWorker) RunCheckConnection() error {
	t := task.Task
	s, err := task.SendToComputer("info", "", "GET")
	if err != nil {
		task.saveWrongConnectionStatus(t)
		return err
	}
	if logLevel {
		dvlog.PrintfFullOnly("Connection info %s", s)
	}
	if t.ConnectionStatus != 0 {
		t.ConnectionStatus = 0
		err = task.saveConnectionStatus(t)
		return err
	}
	return nil
}

func (task *TaskWorker) RunConfigSending() error {
	config := task.Task.Config
	if config == nil {
		return errors.New("No config in task")
	}
	body, err := json.Marshal(config)
	if err != nil {
		return err
	}
	res, err := task.SendToComputer(configUrl, string(body), configMethod)
	if err != nil {
		task.saveWrongConnectionStatus(task.Task)
		return err
	}
	if logLevel {
		dvlog.Print("received from config " + task.Task.Id + " : " + res)
	}
	t:=task.Task
	err = analyzeComputerConfigSendingResponse(res, t)
	if err!=nil {
		return err
	}
	t.ConnectionStatus = 0
	t.TaskStatus = 1
	err = task.saveConfigSending(t)
	return err
}

func (task *TaskWorker) RunFileSending() error {
	fileUrl, body, hint, err := analyzeComputerFileSendingRequest()
	if err != nil {
		return err
	}
	if len(body)==0 {
		t:=task.Task
		t.LeftFiles = nil
		t.ConnectionStatus = 0
		t.TaskStatus = 1000
		err = task.saveFileSending(t) 
		return err
	}
	res, err := task.SendToComputer(fileUrl, string(body), fileSendMethod)
	if err != nil {
		task.saveWrongConnectionStatus(task.Task)
		return err
	}
	if logLevel {
		dvlog.Print("received from file sending " + task.Task.Id + " : " + res)
	}
	t:=task.Task
	err = analyzeComputerFileSendingResponse(res, hint, t)
	if err!=nil {
		return err
	}
	if len(t.LeftFiles)==0 {
		t.LeftFiles = nil 
		t.TaskStatus = 1000
	}
	t.ConnectionStatus = 0
	err = task.saveFileSending(t)
	return err
}

}

func (task *TaskWorker) GetComputerUrl() string {
	s := task.Task.Url
	if s == "" {
		return s
	}
	if s[len(s)-1] != '/' {
		s = s + "/"
	}
	p := strings.Index(s, "//")
	if p < 0 {
		return "http://" + s
	}
	if p == 0 {
		return "http:" + s
	}
	return s
}

func (task *TaskWorker) SendToComputer(url string, body string, method string) (string, error) {
	fullUrl := task.GetComputerUrl() + url
	bodyIo := io.NopCloser(bytes.NewReader([]byte(body)))

	req, err := http.NewRequest(method, fullUrl, bodyIo)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode >= 300 {
		return "", errors.New(strconv.Itoa(res.StatusCode) + " " + url + " " + method + " " + body)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(resBody), nil
}

func (task *TaskWorker) saveWrongConnectionStatus(t *TvTask) error {
	if t.ConnectionStatus < 0 {
		t.ConnectionStatus = 1
	} else {
		t.ConnectionStatus++
	}
	return task.saveConnectionStatus(t)
}

func (task *TaskWorker) saveConnectionStatus(t *TvTask) error {
	newTask, err := createOrUpdateTaskDatabaseForConnectionStatus(t)
	if err != nil {
		return err
	}
	task.Task = newTask
	return nil
}

func (task *TaskWorker) saveConfigSending(t *TvTask) error {
	newTask, err := createOrUpdateTaskDatabaseForConfigSending(t)
	if err != nil {
		return err
	}
	task.Task = newTask
	return nil
}

func (task *TaskWorker) saveFileSending(t *TvTask) error {
	newTask, err := createOrUpdateTaskDatabaseForFileSending(t)
	if err != nil {
		return err
	}
	task.Task = newTask
	return nil
}
