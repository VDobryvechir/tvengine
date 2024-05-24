/***********************************************************************
TV Controller
Copyright 2024 by Volodymyr Dobryvechir (vdobryvechir@gmail.com)
************************************************************************/

package tvcontrol

import (
	"errors"
    "strconv" 
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
)

func readScreens(presentation *dvevaluation.DvVariable) ([]*TvScreen, error) {
	subItem := presentation.ReadSimpleChild("screens")
	if subItem == nil || subItem.Kind != dvevaluation.FIELD_ARRAY || len(subItem.Fields)==0 {
		return nil, errors.New("No screens")
	}
	n:=len(subItem.Fields)
	res:=make([]*TvScreen, n)
	for i:=0 ; i<n ; i++ {
		tv:=&TvScreen
		err = subItem.Fields[i].DvVariableToAnyStruct(tv)
		if err!=nil {
			return nil, err
		}
		res[i] = tv 
	}
	return res, nil
}

func generateConfig(duration []int, screens []*TvScreen) (*TvConfig, error) {
	n:=len(duration)
	m:=len(screens)
	if n==0 || n!=m {
		return nil, errors.New("Strange situation of "+strconv.Itoa(n)+" durations and "+ strconv.Itoa(m)+" screens")
	}
	fl:=make([]string, n)
	for i:=0 ; i<n ; i++ {
		fl[i] = screens[i].FileName
	}
	return &TvConfig{File:fl,Duration:duration}, nil
}

func putUpRealFiles(screens []*TvScreen) ([]string, error) {
	n:=len(screens)
	if n==0 {
		return nil, errors.New("Strange situation of no screens")
	}
	fl:=make([]string, n)
	for i:=0 ; i<n ; i++ {
		fl[i] = screens[i].FileReal
	}
	return fl, nil
}

func prepareSampleTask(presentation *dvevaluation.DvVariable) (*TvTask, error) {
	presId:=presentation.ReadSimpleChildValue("id")
	presName:=presentation.ReadSimpleChildValue("name")
	presVersion:=presentation.ReadSimpleChildValue("version")
	if presId=="" || presName=="" || presVersion=="" {
		return nil, errors.New("id name version must not be empty in presentation "+ id+","+name+","+version)
	}
	duration, err:= presentation.ReadChildIntArrayValue("duration")
	if err!=nil {
		return nil, err
	}
	screens, err:=readScreens(presentation)
	if err!=nil {
		return nil, err
	}
	config, err:=generateConfig(duration, screens)
	if err!=nil {
		return nil, err
	}
	realFiles, err:=putUpRealFiles(screens)
	if err!=nil {
		return nil, err
	}
	r:=&TvTask{NewPresentationId:presId, NewPresentationName: presName, NewPresentationVersion: presVersion, Config: config, RealFiles: realFiles }
	return r, nil
} 

func createTvTasks(sample *TvTask, tvs []*dvevaluation.DvVariable) ([]*TvTask, error) {
	n:=len(tvs)
	if n==0 {
		return nil, errors.New("No tvs")
	}
	res:=make([]*TvTask, n)
	for i:=0; i<n ; i++ {
		tv:=tvs[i]
		id:=tv.ReadSimpleChildValue("id")
		name:=tv.ReadSimpleChildValue("name")
		url:=tv.ReadSimpleChildValue("url")
		if id=="" || name=="" || url=="" {
			return nil, errors.New("empty id, name, url in tvpc "+id+","+name+","+url)
		}
		res[i]= &TvTask{NewPresentationId: sample.NewPresentationId, NewPresentationName: sample.NewPresentationName, NewPresentationVersion: sample.NewPresentationVersion, Config: sample.Config, RealFiles: sample.RealFiles, Id:id, Name:name, Url:url, LeftFiles:make([]string,0,16),ConnectionStatus: -1}
	}
    return res, nil
}

func convertTasksToDvVariables(tasks []*TvTask) (res *dvevaluation.DvVariable,err error) {
	n:=len(tasks)
	res = &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields:make([]*dvevaluation.DvVariable, n)}
	for i:=0 ; i<n ; i++ {
		res.Fields[i], err = dvevaluation.AnyStructToDvVariable(tasks[i])
		if err!=nil {
			return
		}
	}
	return
}
