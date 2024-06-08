/***********************************************************************
TV Controller
Copyright 2024 by Volodymyr Dobryvechir (vdobryvechir@gmail.com)
************************************************************************/

package tvcontrol

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
)

func readScreens(presentation *dvevaluation.DvVariable) ([]*TvScreen, error) {
	subItem := presentation.ReadSimpleChild("screens")
	if subItem == nil || subItem.Kind != dvevaluation.FIELD_ARRAY || len(subItem.Fields) == 0 {
		return nil, errors.New("no screens")
	}
	n := len(subItem.Fields)
	res := make([]*TvScreen, n)
	for i := 0; i < n; i++ {
		tv := &TvScreen{}
		err := subItem.Fields[i].DvVariableToAnyStruct(tv)
		if err != nil {
			return nil, err
		}
		res[i] = tv
	}
	return res, nil
}

func generateConfig(duration []int, screens []*TvScreen) (*TvConfig, error) {
	n := len(duration)
	m := len(screens)
	if n == 0 || n != m {
		return nil, errors.New("Strange situation of " + strconv.Itoa(n) + " durations and " + strconv.Itoa(m) + " screens")
	}
	fl := make([]string, n)
	for i := 0; i < n; i++ {
		fl[i] = screens[i].FileName
	}
	return &TvConfig{File: fl, Duration: duration}, nil
}

func putUpRealFiles(screens []*TvScreen) ([]string, error) {
	n := len(screens)
	if n == 0 {
		return nil, errors.New("strange situation of no screens")
	}
	fl := make([]string, n)
	for i := 0; i < n; i++ {
		fl[i] = screens[i].FileReal
	}
	return fl, nil
}

func prepareSampleTask(presentation *dvevaluation.DvVariable) (*TvTask, error) {
	presId := presentation.ReadSimpleChildValue("id")
	presName := presentation.ReadSimpleChildValue("name")
	presVersion := presentation.ReadSimpleChildValue("version")
	if presId == "" || presName == "" || presVersion == "" {
		return nil, errors.New("id name version must not be empty in presentation " + presId + "," + presName + "," + presVersion)
	}
	duration, err := presentation.ReadChildIntArrayValue("duration")
	if err != nil {
		return nil, err
	}
	screens, err := readScreens(presentation)
	if err != nil {
		return nil, err
	}
	config, err := generateConfig(duration, screens)
	if err != nil {
		return nil, err
	}
	realFiles, err := putUpRealFiles(screens)
	if err != nil {
		return nil, err
	}
	err = fixConfigFileNames(config, realFiles)
	if err != nil {
		return nil, err
	}
	r := &TvTask{NewPresentationId: presId, NewPresentationName: presName, NewPresentationVersion: presVersion, Config: config, RealFiles: realFiles}
	return r, nil
}

func fixConfigFileNames(config *TvConfig, realFiles []string) error {
	n := len(config.File)
	m := len(realFiles)
	if n == 0 {
		return errors.New("no files to show")
	}
	if n != m {
		return errors.New("misconfiguration in screens, remove them and create from the scratch")
	}
	for i := 0; i < n; i++ {
		if len(config.File[i]) == 0 {
			s, err := fixFileName(realFiles[i])
			if err != nil {
				return err
			}
			config.File[i] = s
			dvlog.PrintfFullOnly("Fixed %d file name %s", i, s)
		}
	}
	return nil
}

func getFileNamePrefix(ext string) (string, error) {
	if ext == "mp4" || ext == "webm" || ext == "ogv" {
		return "v", nil
	}
	if ext == "jpg" || ext == "jpeg" || ext == "png" || ext == "gif" {
		return "i", nil
	}
	return "", errors.New("unsupported file format " + ext)
}

func getFileNameExtension(name string) (string, error) {
	p := strings.LastIndex(name, ".")
	if p <= 0 {
		return "", errors.New("File " + name + " must have extension")
	}
	ext := strings.ToLower(name[p+1:])
	return ext, nil
}

func getFileNameLength(name string) (string, error) {
	file := dvparser.GetByGlobalPropertiesOrDefault("HTML_PATH", "") + name
	fi, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	size := fi.Size()
	if size == 0 {
		return "", errors.New("File " + file + " has zero size")
	}
	res := strconv.FormatInt(size, 10)
	return res, nil
}

func getFileNameId(name string) (string, error) {
	begPos := 0
	n := len(name)
	for begPos < n && !(name[begPos] >= '0' && name[begPos] <= '9') {
		begPos++
	}
	if begPos == n {
		return "", errors.New("No id in file name " + name)
	}
	endPos := begPos + 1
	for endPos < n && name[endPos] >= '0' && name[endPos] <= '9' {
		endPos++
	}
	return name[begPos:endPos], nil
}

// template: i583747_7721532218530737715-1071263.png
func fixFileName(name string) (string, error) {
	resLen, err := getFileNameLength(name)
	if err != nil {
		return "", err
	}
	resExt, err := getFileNameExtension(name)
	if err != nil {
		return "", err
	}
	resPref, err := getFileNamePrefix(resExt)
	if err != nil {
		return "", err
	}
	resId, err := getFileNameId(name)
	if err != nil {
		return "", err
	}
	return resPref + resId + "_0-" + resLen + "." + resExt, nil
}

func createTvTasks(sample *TvTask, tvs []*dvevaluation.DvVariable) ([]*TvTask, error) {
	n := len(tvs)
	if n == 0 {
		return nil, errors.New("no tvs")
	}
	res := make([]*TvTask, n)
	for i := 0; i < n; i++ {
		tv := tvs[i]
		id := tv.ReadSimpleChildValue("id")
		name := tv.ReadSimpleChildValue("name")
		url := tv.ReadSimpleChildValue("url")
		if id == "" || name == "" || url == "" {
			return nil, errors.New("empty id, name, url in tvpc " + id + "," + name + "," + url)
		}
		res[i] = &TvTask{NewPresentationId: sample.NewPresentationId, NewPresentationName: sample.NewPresentationName, NewPresentationVersion: sample.NewPresentationVersion, Config: sample.Config, RealFiles: sample.RealFiles, Id: id, Name: name, Url: url, LeftFiles: make([]string, 0, 16), ConnectionStatus: -1}
	}
	return res, nil
}
