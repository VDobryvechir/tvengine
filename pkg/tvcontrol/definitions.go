/***********************************************************************
TV Controller
Copyright 2024 by Volodymyr Dobryvechir (vdobryvechir@gmail.com)
************************************************************************/

package tvcontrol

type TvConfig struct {
	File []string `json:"file"`
	Duration []int `json:"duration"`
}

type TvScreen struct {
	FileReal string `json:"file"`
	FileName string `json:"fileName"`
	Id       string `json:"id"`
}

type TvTask struct {
	Id  string `json:"id"`,
	Name string `json:"name"`
	Url  string `json:"url"`
	OldPresentationId string `json:"oldPresentationId"`
	OldPresentationName string `json:"oldPresentationName"`
	OldPresentationVersion string `json:"oldPresentationVersion"`
	NewPresentationId string `json:"newPresentationId"`
	NewPresentationName string `json:"newPresentationName"`
	NewPresentationVersion string `json:"newPresentationVersion"`
	Config *TvConfig `json:"config"`
	RealFiles []string `json:"realFiles"`
	LeftFiles []string `json:"leftFiles"`
	TaskStatus int `json:"taskStatus"`
	ConnectionStatus int `json:"connectionStatus"`
}