package cvt2json

import (
	"log"
	"os"
	"testing"
)

func TestConvert(t *testing.T) {
	data, err := os.ReadFile("./StudentPersonals.json")
	if err != nil {
		log.Fatalln(err)
	}
	Cvt2Pesc(string(data), "./temp.json")
}
