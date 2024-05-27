/***********************************************************************
TV Controller
Copyright 2024 by Volodymyr Dobryvechir (vdobryvechir@gmail.com)
************************************************************************/

package tvcontrol

import (
	"github.com/Dobryvechir/microcore/pkg/dvdbmanager"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
)

const taskDbName = "task"

var taskConditionsForWeb = []string{
	"NEW",
	"previous.NewPresentationVersion == current.NewPresentationVersion && previous.Url == current.Url && previous.NewPresentationId == current.NewPresentationId",
	"DEFAULT",
}

var taskFieldsForWeb = []string{
	"",
	"OldPresentationId,OldPresentationName,OldPresentationVersion,LeftFiles,TaskStatus,ConnectionStatus",
	"OldPresentationId,OldPresentationName,OldPresentationVersion,ConnectionStatus",
}

const taskConditionsForConfigSendingPart1 = "current.NewPresentationVersion=="
const taskConditionsForConfigSendingPart2 = " && current.NewPresentationId=="

var taskFieldsForConfigSending = []string{
	"!OldPresentationId,OldPresentationName,OldPresentationVersion",
	"Name,NewPresentationName",
	"Name,NewPresentationId,NewPresentationName,NewPresentationVersion,Config,RealFiles,LeftFiles,TaskStatus",
}

var taskFieldsForFileSending = []string{
	"!OldPresentationId,OldPresentationName,OldPresentationVersion",
	"Name,NewPresentationName",
	"^OldPresentationId,OldPresentationName,OldPresentationVersion,ConnectionStatus",
}

var taskConditionsForConnectionCheck = []string{
	"DEFAULT",
}

// all fields except ConnectionStatus must be here
var taskFieldsForConnectionCheck = []string{
	"^ConnectionStatus",
}

func createOrUpdateTaskDatabaseForWeb(tasks []*TvTask) (res []*dvevaluation.DvVariable, err error) {
	n := len(tasks)
	res = make([]*dvevaluation.DvVariable, n)
	for i := 0; i < n; i++ {
		res[i], err = createOrUpdateTaskDatabase(tasks[i], taskConditionsForWeb, taskFieldsForWeb)
		if err != nil {
			return
		}
	}
	return
}

func createOrUpdateTaskDatabase(task *TvTask, taskConditions []string, taskFields []string) (*dvevaluation.DvVariable, error) {
	rowTask, err := dvevaluation.AnyStructToDvVariable(task)
	if err != nil {
		return nil, err
	}
	res, err := dvdbmanager.CreateOrUpdateByConditionsAndUpdateFields(taskDbName, rowTask, taskConditions, taskFields)
	return res, err
}

func getCoincidenceInTask(task *TvTask) string {
	return taskConditionsForConfigSendingPart1 + task.NewPresentationVersion + taskConditionsForConfigSendingPart2 + task.NewPresentationVersion
}

func getCoincidenceConditions(task *TvTask) []string {
	return []string{"current.Url != previous.Url", getCoincidenceInTask(task), "DEFAULT"}
}

func createOrUpdateTaskDatabaseForConfigSending(task *TvTask) (*TvTask, error) {
	rowTask, err := dvevaluation.AnyStructToDvVariable(task)
	if err != nil {
		return nil, err
	}
	task.OldPresentationId = task.NewPresentationId
	task.OldPresentationName = task.NewPresentationName
	task.OldPresentationVersion = task.NewPresentationVersion
	taskConditions := getCoincidenceConditions(task)
	res, err := dvdbmanager.CreateOrUpdateByConditionsAndUpdateFields(taskDbName, rowTask, taskConditions, taskFieldsForConfigSending)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	tsk := &TvTask{}
	err = res.DvVariableToAnyStruct(tsk)
	return tsk, err
}

// it is assumed that the only changed fiedls are LeftFiles and TaskStatus, possibly also ConnectionStatus
func createOrUpdateTaskDatabaseForFileSending(task *TvTask) (*TvTask, error) {
	rowTask, err := dvevaluation.AnyStructToDvVariable(task)
	if err != nil {
		return nil, err
	}
	taskConditions := getCoincidenceConditions(task)
	res, err := dvdbmanager.CreateOrUpdateByConditionsAndUpdateFields(taskDbName, rowTask, taskConditions, taskFieldsForConfigSending)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	tsk := &TvTask{}
	err = res.DvVariableToAnyStruct(tsk)
	return tsk, err
}

// it is assumed that the only changed fiedls is ConnectionStatus
func createOrUpdateTaskDatabaseForConnectionStatus(task *TvTask) (*TvTask, error) {
	rowTask, err := dvevaluation.AnyStructToDvVariable(task)
	if err != nil {
		return nil, err
	}
	res, err := dvdbmanager.CreateOrUpdateByConditionsAndUpdateFields(taskDbName, rowTask, taskConditionsForConnectionCheck, taskFieldsForConnectionCheck)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	tsk := &TvTask{}
	err = res.DvVariableToAnyStruct(tsk)
	return tsk, err
}
