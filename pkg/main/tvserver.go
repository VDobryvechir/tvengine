// package main is the entry point

package main

import (
	_ "github.com/Dobryvechir/microcore/pkg/dvoc"
    _ "github.com/VDobryvechir/tvengine/pkg/tvcontrol"
	"github.com/Dobryvechir/microcore/pkg/dvconfig"
)

func main() {
	dvconfig.SetApplicationName("tvserver")
	dvconfig.ServerStart()
}
