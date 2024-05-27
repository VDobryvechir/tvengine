/***********************************************************************
TV Controller
Copyright 2024 by Volodymyr Dobryvechir (vdobryvechir@gmail.com)
************************************************************************/

package tvcontrol

import (
	"time"

	"github.com/Dobryvechir/microcore/pkg/dvdbmanager"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
)

var taskWorkerPoolChannel = make(chan int)
var taskWorkerPool map[string]*TaskWorker = make(map[string]*TaskWorker)
var logLevel = true

func wakeUpMainWorker() error {
	select {
	case taskWorkerPoolChannel <- 0:
		if logLevel {
			dvlog.Print("woke up main worker")
		}
	default:
		if logLevel {
			dvlog.Print("main worker already waken up")
		}
	}
	return nil
}

func RunMainWorker() {
	time.Sleep(5 * time.Second)
	for {
		res, err := dvdbmanager.RecordReadAll(taskDbName)
		if err == nil {
			loadMainTasks(res)
		} else {
			dvlog.PrintError(err)
		}
		<-taskWorkerPoolChannel
	}

}

func loadMainTasks(res *dvevaluation.DvVariable) {
	if res == nil || len(res.Fields) == 0 {
		dvlog.PrintlnError("No task is defined yet")
		return
	}
	rest := make(map[string]int)
	n := len(res.Fields)
	for i := 0; i < n; i++ {
		v := res.Fields[i]
		id := v.ReadSimpleChildValue("id")
		if len(id) == 0 {
			continue
		}
		rest[id] = 1
		createOrWakeUpTaskById(id, v)
	}
	for k := range taskWorkerPool {
		if rest[k] != 1 {
			getTaskDownById(k)
		}
	}
}

func getTaskDownById(id string) {
	task, ok := taskWorkerPool[id]
	if !ok {
		return
	}
	delete(taskWorkerPool, id)
	if task == nil {
		return
	}
	select {
	case task.StopChannel <- 2:
		if logLevel {
			dvlog.Print("Stop sent")
		}
	default:
		if logLevel {
			dvlog.Print("Stop already sent")
		}
	}
}

func createOrWakeUpTaskById(id string, v *dvevaluation.DvVariable) {
	task, ok := taskWorkerPool[id]
	if !ok || task == nil {
		tvTask := &TvTask{}
		err := v.DvVariableToAnyStruct(tvTask)
		if err != nil {
			tvTask = nil
			dvlog.PrintError(err)
		}
		task = &TaskWorker{Id: id, Task: tvTask, WakeUpChannel: make(chan int), StopChannel: make(chan int)}
		taskWorkerPool[id] = task
		go func() {
			task.RunBackground()
			dvlog.Print("Task " + task.Id + " ended")
			if _, ok = taskWorkerPool[id]; ok {
				taskWorkerPool[id] = nil
			}
		}()
	} else {
		select {
		case task.WakeUpChannel <- 1:
			if logLevel {
				dvlog.Print("Wake up sent")
			}
		default:
			if logLevel {
				dvlog.Print("Wake up already sent")
			}
		}
	}
}
