/***********************************************************************
TV Controller
Copyright 2024 by Volodymyr Dobryvechir (vdobryvechir@gmail.com)
************************************************************************/

package tvcontrol

import (
	"github.com/Dobryvechir/microcore/pkg/dvaction"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	_ "github.com/Dobryvechir/microcore/pkg/dvdbmanager"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"

)

type TvControlConfig struct {
	Presentation  string `json:"presentation"`
	Tv   string `json:"tv"`
	Result   string `json:"result"`
}

func TvControlInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &TvControlConfig{}
	if !dvaction.DefaultInitWithObject(command, config, dvaction.GetEnvironment(ctx)) {
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func TvControlRun(data []interface{}) bool {
	config := data[0].(*TvControlConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return tvControlRunByConfig(config, ctx)
}

func tvControlRunByConfig(config *TvControlConfig, ctx *dvcontext.RequestContext) bool {
	presentationData, ok := dvaction.ReadActionResult(config.Presentation, ctx)
	if !ok {
		dvlog.PrintlnError("Presentation is absent")
		return true
	}
	presentation := dvevaluation.AnyToDvVariable(presentationData)
	if presentation == nil || presentation.Kind != dvevaluation.FIELD_OBJECT || len(presentation.Fields) == 0 {
		dvlog.PrintlnError("Presentation is empty")
		return true
	}
	tvData, ok := dvaction.ReadActionResult(config.Tv, ctx)
	if !ok {
		dvlog.PrintlnError("tv data is absent")
		return true
	}
	tv := dvevaluation.AnyToDvVariable(tvData)
	if tv == nil || tv.Kind != dvevaluation.FIELD_OBJECT || len(tv.Fields) == 0 {
		dvlog.PrintlnError("tv data is empty")
		return true
	}
	n:=len(tv.Fields)
	res:=&dvevaluation.DvVariable{Kind:dvevaluation.FIELD_ARRAY, Fields: make([]*dvevaluation.DvVariable, 0, n)}
	sample, err:=prepareSampleTask(presentation)
	if err!=nil {
		dvlog.PrintError(err)
		return true
	}
	tasks,err:=createTvTasks(sample, tv)
	if err!=nil {
		dvlog.PrintError(err)
		return true
	}
	res, err = createOrUpdateTaskDatabaseForWeb(tasks)
	if err!=nil {
		dvlog.PrintError(err)
		return true
	}
	err = wakeUpMainWorker()
	if err!=nil {
		dvlog.PrintError(err)
		return true
	}
	dvaction.SaveActionResult(config.Result, res, ctx)
	return true
}

const (
	CommandTvControl     = "tvcontrol"
)

var processFunctions = map[string]dvaction.ProcessFunction{
	CommandTvControl:     {Init: TvControlInit, Run: TvControlRun},
}

func Init() bool {
	dvaction.AddProcessFunctions(processFunctions)
	return true
}

var inited = Init()

