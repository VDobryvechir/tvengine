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
	"previous.newPresentationVersion == current.newPresentationVersion && previous.url == current.url && previous.newPresentationId == current.newPresentationId",
	"DEFAULT",
}

var taskFieldsForWeb = []string{
	"",
	"oldPresentationId,oldPresentationName,oldPresentationVersion,leftFiles,taskStatus,connectionStatus",
	"oldPresentationId,oldPresentationName,oldPresentationVersion,connectionStatus",
}

const taskConditionsForConfigSendingPart1 = "current.newPresentationVersion=="
const taskConditionsForConfigSendingPart2 = " && current.newPresentationId=="

var taskFieldsForConfigSending = []string{
	"!oldPresentationId,oldPresentationName,oldPresentationVersion",
	"name,newPresentationName",
	"name,newPresentationId,newPresentationName,newPresentationVersion,config,realFiles,leftFiles,taskStatus",
}

var taskFieldsForFileSending = []string{
	"!oldPresentationId,oldPresentationName,oldPresentationVersion",
	"name,newPresentationName",
	"^oldPresentationId,oldPresentationName,oldPresentationVersion,connectionStatus",
}

var taskConditionsForConnectionCheck = []string{
	"DEFAULT",
}

// all fields except ConnectionStatus must be here
var taskFieldsForConnectionCheck = []string{
	"^connectionStatus",
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
	return taskConditionsForConfigSendingPart1 + task.NewPresentationVersion + taskConditionsForConfigSendingPart2 + task.NewPresentationId
}

func getCoincidenceConditions(task *TvTask) []string {
	return []string{"current.url != previous.url", getCoincidenceInTask(task), "DEFAULT"}
}

func createOrUpdateTaskDatabaseForConfigSending(task *TvTask) (*TvTask, error) {
	task.OldPresentationId = task.NewPresentationId
	task.OldPresentationName = task.NewPresentationName
	task.OldPresentationVersion = task.NewPresentationVersion
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
