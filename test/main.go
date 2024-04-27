// package main is the entry point

package main

import (
	"github.com/Dobryvechir/microcore/pkg/dvconfig"
	_ "github.com/Dobryvechir/microcore/pkg/dvdbdata"
	_ "github.com/Dobryvechir/microcore/pkg/dvzoo"
	_ "github.com/Dobryvechir/microcore/pkg/dvoc"
	_ "github.com/lib/pq"
)

func main() {
	dvconfig.SetApplicationName("tvserver")
	dvconfig.ServerStart()
}
