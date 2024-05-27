/***********************************************************************
TV Controller
Copyright 2024 by Volodymyr Dobryvechir (vdobryvechir@gmail.com)
************************************************************************/

package tvcontrol

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
)

const packageSize = 1 << 19

const configUrl = "config"
const configMethod = "POST"

const fileSendUrl = "upload/"
const fileSendMethod = "POST"

func getLeftFiles(r string) ([]string, error) {
	res := make(map[string]int, 10)
	err := json.Unmarshal([]byte(r), res)
	if err != nil {
		return nil, err
	}
	d := make([]string, 0, 16)
	for k, v := range res {
		s := k
		if v != 0 {
			s += ":" + strconv.Itoa(v)
		}
		d = append(d, s)
	}
	return d, nil
}

func analyzeComputerConfigSendingResponse(r string, t *TvTask) error {
	leftFiles, err := getLeftFiles(r)
	if err != nil {
		return err
	}
	t.LeftFiles = leftFiles
	return nil
}

func analyzeComputerFileSendingResponse(r string, hint string, t *TvTask) error {
	if len(hint) == 0 {
		t.LeftFiles = t.LeftFiles[1:]
	} else if len(t.LeftFiles) > 0 {
		t.LeftFiles[0] = hint
	}
	if logLevel {
		dvlog.PrintfFullOnly("Left %v, response %s", t.LeftFiles, r)
	}
	m := len(t.RealFiles)
	p := len(t.LeftFiles)
	if m == 0 || p == 0 {
		t.TaskStatus = 1000
		t.LeftFiles = nil
	} else {
		done := (m-p)*999/m + 1
		subCurrent := getCurrentSeek(t.LeftFiles[0])
		subTotal := getFullSeek(t.LeftFiles[0])
		if subTotal > 0 {
			done += subCurrent * 999 / (m * subTotal)
		}
		if done > 999 {
			done = 999
		}
		t.TaskStatus = done
	}
	return nil
}

func getCurrentSeek(s string) int {
	pos := strings.Index(s, ":")
	if pos > 0 {
		t := s[pos+1:]
		n, err := strconv.Atoi(t)
		if err != nil {
			dvlog.PrintlnError("strange info at " + s + " " + err.Error())
			n = 0
		}
		return n
	}
	return 0
}

func getFullSeek(s string) int {
	pos := strings.Index(s, ":")
	if pos > 0 {
		s = s[:pos]
	}
	pos = strings.LastIndex(s, "-")
	if pos > 0 {
		t := s[pos+1:]
		n, err := strconv.Atoi(t)
		if err != nil {
			dvlog.PrintlnError("strange file size at " + s + " " + err.Error())
			n = 0
		}
		return n
	}
	return 0
}

func changeSeek(s string, newSeek int) string {
	pos := strings.Index(s, ":")
	if pos > 0 {
		s = s[:pos]
	}
	s += ":" + strconv.Itoa(newSeek)
	return s
}

func analyzeComputerFileSendingRequest(t *TvTask) (url string, body string, hint string, err error) {
	n := len(t.LeftFiles)
	for n > 0 {
		p := t.LeftFiles[0]
		current := getCurrentSeek(p)
		total := getFullSeek(p)
		if current >= total {
			t.LeftFiles = t.LeftFiles[1:]
			continue
		}
		dif := total - current
		if dif > packageSize {
			dif = packageSize
			hint = changeSeek(t.LeftFiles[0], current+dif)
		}
		index, name, err2 := detectRealFileName(t, t.LeftFiles[0])
		if err2 != nil {
			dvlog.PrintError(err2)
			t.LeftFiles = t.LeftFiles[1:]
			continue
		}
		url = fileSendUrl + strconv.Itoa(index) + "_" + strconv.Itoa(current)
		data, err1 := readFileWithSeek(name, current, dif)
		if err1 != nil {
			err = err1
			return
		}
		body = string(data)
		return
	}
	return
}

func readFileWithSeek(name string, seek int, amount int) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := make([]byte, amount)
	_, err = f.Seek(int64(seek), io.SeekStart)
	if err != nil {
		return nil, err
	}
	n, err := f.Read(buf)
	if err != nil {
		return nil, err
	}
	if n != amount {
		return nil, errors.New("File " + name + " has insufficient size after " + strconv.Itoa(seek) + " for " + strconv.Itoa(amount) + " but " + strconv.Itoa(n))
	}
	return buf, nil
}

func detectRealFileName(t *TvTask, s string) (index int, file string, err error) {
	pos := strings.Index(s, ":")
	if pos > 0 {
		s = s[:pos]
	}
	index = dvtextutils.FindIndexInStringArray(t.Config.File, s)
	if index < 0 {
		return index, "", errors.New("file " + s + " is not detected in config")
	}
	if index >= len(t.RealFiles) {
		return index, "", errors.New("file " + s + " is not provided in real-file list")
	}
	file = dvparser.GetByGlobalPropertiesOrDefault("HTML_PATH", "") + t.RealFiles[index]
	return
}
