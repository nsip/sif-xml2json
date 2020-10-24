package main

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	cfg "github.com/nsip/sif-xml2json/config/cfg"
)

func TestMain(t *testing.T) {
	main()
}

func TestLoad(t *testing.T) {
	c := cfg.NewCfg(
		"Config",
		map[string]string{
			"[s]": "Service",
			"[v]": "Version",
		},
		"../config/config.toml",
	).(*cfg.Config)
	spew.Dump(*c)
}

func TestInit(t *testing.T) {
	c := cfg.NewCfg(
		"Config",
		map[string]string{
			"[s]":    "Service",
			"[v]":    "Version",
			"[port]": "WebService.Port",
		},
		"../config/config.toml",
	).(*cfg.Config)
	spew.Dump(*c)

	c = env2Struct("Config", &cfg.Config{}).(*cfg.Config)
	spew.Dump(*c)
}
