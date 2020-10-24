package main

import (
	"fmt"
	"os"

	"github.com/cdutwhu/gonfig/strugen"
)

func main() {
	pkgName := "cfg"
	cfgSrc := fmt.Sprintf("./%s/config.go", pkgName)
	os.Remove(cfgSrc)
	strugen.GenStruct("./config.toml", "Config", pkgName, cfgSrc)
	strugen.GenNewCfg(cfgSrc)
}
