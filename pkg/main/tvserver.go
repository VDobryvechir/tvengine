// package main is the entry point

package main

import (
	_ "github.com/Dobryvechir/microcore/pkg/dvoc"
    "github.com/VDobryvechir/tvengine/pkg/tvcontrol"
    "github.com/Dobryvechir/microcore/pkg/dvconfig"
)

func main() {
	dvconfig.SetApplicationName("tvserver")
        tvcontrol.RunMainWorker()
	dvconfig.ServerStart()
}
