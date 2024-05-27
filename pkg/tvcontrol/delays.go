/***********************************************************************
TV Controller
Copyright 2024 by Volodymyr Dobryvechir (vdobryvechir@gmail.com)
************************************************************************/

package tvcontrol

var delayInErrorCase = 30
var delayInIdleCase = 60
var delayInOperationCase = 0

func GetDelayInErrorCase() int {
	return delayInErrorCase
}

func GetDelayInIdleCase() int {
	return delayInIdleCase
}

func GetDelayInOperationCase() int {
	return delayInOperationCase
}
