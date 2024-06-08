/***********************************************************************
TV Controller
Copyright 2024 by Volodymyr Dobryvechir (vdobryvechir@gmail.com)
************************************************************************/

package tvcontrol

import (
	"errors"

	"github.com/Dobryvechir/microcore/pkg/dvaction"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	_ "github.com/Dobryvechir/microcore/pkg/dvdbmanager"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
)

type TvControlConfig struct {
	Presentation string `json:"presentation"`
	Tv           string `json:"tv"`
	Result       string `json:"result"`
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
	err := tvControlRunByConfig(config, ctx)
	if err != nil {
		mes := err.Error()
		dvlog.PrintlnError(mes)
		resError := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_STRING, Name: []byte("error"), Value: []byte(mes)}
		res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_OBJECT, Fields: []*dvevaluation.DvVariable{resError}}
		dvaction.SaveActionResult(config.Result, res, ctx)
	}
	return true
}

func tvControlRunByConfig(config *TvControlConfig, ctx *dvcontext.RequestContext) error {
	presentationData, ok := dvaction.ReadActionResult(config.Presentation, ctx)
	if !ok {
		return errors.New("system error in reading the presentation")
	}
	presentation := dvevaluation.AnyToDvVariable(presentationData)
	if presentation == nil || presentation.Kind != dvevaluation.FIELD_OBJECT || len(presentation.Fields) == 0 {
		return errors.New("cannot send empty presentation")
	}
	tvData, ok := dvaction.ReadActionResult(config.Tv, ctx)
	if !ok {
		return errors.New("system error in reading tv data")
	}
	tv := dvevaluation.AnyToDvVariable(tvData)
	if tv == nil || tv.Kind != dvevaluation.FIELD_ARRAY || len(tv.Fields) == 0 {
		return errors.New("there is no tv pc is current group")
	}
	n := len(tv.Fields)
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields: make([]*dvevaluation.DvVariable, 0, n)}
	sample, err := prepareSampleTask(presentation)
	if err != nil {
		return err
	}
	tasks, err := createTvTasks(sample, tv.Fields)
	if err != nil {
		return err
	}
	res.Fields, err = createOrUpdateTaskDatabaseForWeb(tasks)
	if err != nil {
		return err
	}
	err = wakeUpMainWorker()
	if err != nil {
		return err
	}
	dvaction.SaveActionResult(config.Result, res, ctx)
	return nil
}

const (
	CommandTvControl = "tvcontrol"
)

var processFunctions = map[string]dvaction.ProcessFunction{
	CommandTvControl: {Init: TvControlInit, Run: TvControlRun},
}

func Init() bool {
	dvaction.AddProcessFunctions(processFunctions)
	return true
}

var inited = Init()
