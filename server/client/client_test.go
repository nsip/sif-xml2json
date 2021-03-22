package client

import (
	"os"
	"testing"
)

func TestDO(t *testing.T) {
	config := "./config.toml"
	str, err := DO(
		config,
		"Help",
		nil,
	)
	fPln(str)
	fPln(err)
	fPln(" ------------------------------------ ")

	bytes, err := os.ReadFile("../../data/examples/3.4.8/NAPCodeFrame_0.xml")
	failOnErr("%v", err)
	str, err = DO(
		config,
		"Convert",
		&Args{
			Data:   bytes,
			Ver:    "3.4.8",
			ToNATS: false,
		},
	)
	fPln(str)
	fPln(err)
	mustWriteFile("./out.json", []byte(str))
}
