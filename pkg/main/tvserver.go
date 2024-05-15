// package main is the entry point

package main

import (
	"github.com/Dobryvechir/microcore/pkg/dvconfig"
	_ "github.com/Dobryvechir/microcore/pkg/dvoc"
)

func main() {
	dvconfig.SetApplicationName("tvserver")
	dvconfig.ServerStart()
}
