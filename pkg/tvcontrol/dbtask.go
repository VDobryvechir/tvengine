/***********************************************************************
TV Controller
Copyright 2024 by Volodymyr Dobryvechir (vdobryvechir@gmail.com)
************************************************************************/

package tvcontrol

import (
	"strconv"
	"github.com/Dobryvechir/microcore/pkg/dvdbmanager"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
)

const taskDbName = "task"

const taskConditionsForWeb = [
	"NEW",
	"previous.NewPresentationVersion == current.NewPresentationVersion && previous.Url == current.Url && previous.NewPresentationId == current.NewPresentationId"
	"DEFAULT"
]

const taskFieldsForWeb = [
    "",
	"OldPresentationId,OldPresentationName,OldPresentationVersion,LeftFiles,TaskStatus,ConnectionStatus"
	"OldPresentationId,OldPresentationName,OldPresentationVersion,ConnectionStatus",
]


const taskConditionsForConfigSendingPart1 = "current.NewPresentationVersion=="
const taskConditionsForConfigSendingPart2 = " && current.NewPresentationId=="

const taskFieldsForConfigSending = [
	"!OldPresentationId,OldPresentationName,OldPresentationVersion",
	"Name,NewPresentationName",
	"Name,NewPresentationId,NewPresentationName,NewPresentationVersion,Config,RealFiles,LeftFiles,TaskStatus",
]
func createOrUpdateTaskDatabaseForWeb(tasks []*TvTask) (res []*dvevaluation.DvVariable,err error) {
	n:=len(tasks)
	res = make([]*dvevaluation.DvVariable, n)
	for i:=0 ; i<n ; i++ {
		res[i], err=createOrUpdateTaskDatabase(tasks[i], taskConditionsForWeb, taskFieldsForWeb)
		if err!=nil {
			return
		}
	}
	return
}

func createOrUpdateTaskDatabase(task *TvTask, taskConditions []string, taskFields []string) (*dvevaluation.DvVariable, error) {
	rowTask, err := dvevaluation.AnyStructToDvVariable(task)
    if err!=nil {
		return nil, err
	}
	res, err:= dvdbmanager.CreateOrUpdateByConditionsAndUpdateFields(taskDbName, rowTask, taskConditions, taskFields)
	return res, err
}

func getCoincidenceInTask(task *TvTask) string {
	return taskConditionsForConfigSendingPart1 + strconv.Itoa(task.NewPresentationVersion) + taskConditionsForConfigSendingPart2 + strconv.Itoa(task.NewPresentationVersion)
}

func getCoincidenceConditions(task *TvTask) []string {
	return []string{"current.Url != previous.Url",getCoincidenceInTask(task), "DEFAULT"}
}

func createOrUpdateTaskDatabaseForConfigSending(task *TvTask) (task *TvTask, error) {
	task.OldPresentationId = task.NewPresentationId
	task.OldPresentationName = task.NewPresentationName
	task.OldPresentationVersion = task.NewPresentationVersion
	taskConditions:=getCoincidenceConditions(task)
	res, err:= dvdbmanager.CreateOrUpdateByConditionsAndUpdateFields(taskDbName, rowTask, taskConditions, taskFieldsForConfigSending)
	if err!=nil {
		return nil, err
	}
	if res==nil {
		return nil, nil
	}
	tsk:=&TvTask{}
	err=res.DvVariableToAnyStruct(tsk)
	return tsk, err
}

