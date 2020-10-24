package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gonfig/strugen"
	"github.com/cdutwhu/gotil/io"
	"github.com/cdutwhu/gotil/rflx"
	"github.com/cdutwhu/gotil/str"
)

var (
	fPln           = fmt.Println
	fPf            = fmt.Printf
	fSf            = fmt.Sprintf
	sHasPrefix     = strings.HasPrefix
	sSplit         = strings.Split
	sReplace       = strings.Replace
	sCount         = strings.Count
	sTrim          = strings.Trim
	mapKeys        = rflx.MapKeys
	rmHeadToFirst  = str.RmHeadToFirst
	rmHeadToLast   = str.RmHeadToLast
	rmTailFromLast = str.RmTailFromLast
	mustWriteFile  = io.MustWriteFile
	failOnErr      = fn.FailOnErr
	failOnErrWhen  = fn.FailOnErrWhen
)

// Println :
func Println(num bool, slc ...string) {
	if num {
		for i, v := range slc {
			fmt.Printf("%d: %v\n", i, v)
		}
	} else {
		for _, v := range slc {
			fmt.Println(v)
		}
	}
}

// ObjGrp :
func ObjGrp(sep string, listGrp ...string) []string {
	m := map[string]bool{}
	for _, lsPath := range listGrp {
		obj := sSplit(lsPath, sep)[0]
		if _, ok := m[obj]; !ok {
			m[obj] = true
		}
	}
	return mapKeys(m).([]string)
}

// MapOfGrp :
func MapOfGrp(objs []string, sep string, xxxPathGrp ...string) map[string][]string {
	m := make(map[string][]string)
	for _, obj := range objs {
		prefix := obj + sep
		for _, lp := range xxxPathGrp {
			if sHasPrefix(lp, prefix) {
				m[obj] = append(m[obj], rmHeadToFirst(lp, sep))
			}
		}
	}
	return m
}

// PrintGrp4Cfg :
func PrintGrp4Cfg(m map[string][]string, attr string) (toml string) {
	for obj, grp := range m {
		content := fmt.Sprintf("[%s]\n  %s = [", obj, attr)
		for _, path := range grp {
			content += fmt.Sprintf("\"%s\", ", path)
		}
		toml += content[:len(content)-2] + "]" + "\n\n"
	}
	return
}

// GenTomlAndGoSrc :
func GenTomlAndGoSrc(specPath, outDir string) {

	const (
		SEP     = "/"
		VERSION = "VERSION: "
		OBJECT  = "OBJECT: "
		LIST    = "LIST: "
		NUMERIC = "NUMERIC: "
		BOOLEAN = "BOOLEAN: "
	)

	var (
		objGrp      []string
		listPathGrp []string
		numPathGrp  []string
		boolPathGrp []string
		SIFVer      string
	)

	bytes, err := ioutil.ReadFile(specPath)
	failOnErr("%v", err)

	for _, line := range sSplit(string(bytes), "\n") {
		switch {
		case sHasPrefix(line, VERSION):
			SIFVer = sTrim(line[len(VERSION):], " \t\r\n")
		case sHasPrefix(line, OBJECT):
			objGrp = append(objGrp, sTrim(line[len(OBJECT):], " \t\r\n"))
		case sHasPrefix(line, LIST):
			// listPathGrp = append(listPathGrp, rmTailFromLast(line[len(LIST):], "/")) // exclude last one
			listPathGrp = append(listPathGrp, sTrim(line[len(LIST):], " \t\r\n"))
		case sHasPrefix(line, NUMERIC):
			numPathGrp = append(numPathGrp, sTrim(line[len(NUMERIC):], " \t\r\n"))
		case sHasPrefix(line, BOOLEAN):
			boolPathGrp = append(boolPathGrp, sTrim(line[len(BOOLEAN):], " \t\r\n"))
		}
	}

	// Println(true, objGrp...)
	// fmt.Println("-----------------------------")

	// Println(false, listPathGrp...)
	// fmt.Println("-----------------------------")

	mListAttr := MapOfGrp(ObjGrp(SEP, listPathGrp...), SEP, listPathGrp...)
	mNumAttr := MapOfGrp(ObjGrp(SEP, numPathGrp...), SEP, numPathGrp...)
	mBoolAttr := MapOfGrp(ObjGrp(SEP, boolPathGrp...), SEP, boolPathGrp...)

	verln := fSf("Version = \"%s\"\n\n", SIFVer)
	toml4List := verln + PrintGrp4Cfg(mListAttr, "LIST")
	toml4Num := verln + PrintGrp4Cfg(mNumAttr, "NUMERIC")
	toml4Bool := verln + PrintGrp4Cfg(mBoolAttr, "BOOLEAN")

	mustWriteFile(outDir+"toml/List2JSON.toml", []byte(toml4List))
	mustWriteFile(outDir+"toml/Num2JSON.toml", []byte(toml4Num))
	mustWriteFile(outDir+"toml/Bool2JSON.toml", []byte(toml4Bool))
}

func main() {

	cfgSrc, pkgName := "./2_toml2json/config.go", "main"
	os.Remove(cfgSrc)

	GenTomlAndGoSrc("./3.4.6.txt", "./3.4.6/")
	toml346 := "./3.4.6/toml/"
	strugen.GenStruct(toml346+"List2JSON.toml", "CfgL2J346", pkgName, cfgSrc)
	strugen.GenStruct(toml346+"Bool2JSON.toml", "CfgB2J346", pkgName, cfgSrc)
	strugen.GenStruct(toml346+"Num2JSON.toml", "CfgN2J346", pkgName, cfgSrc)

	GenTomlAndGoSrc("./3.4.7.txt", "./3.4.7/")
	toml347 := "./3.4.7/toml/"
	strugen.GenStruct(toml347+"List2JSON.toml", "CfgL2J347", pkgName, cfgSrc)
	strugen.GenStruct(toml347+"Bool2JSON.toml", "CfgB2J347", pkgName, cfgSrc)
	strugen.GenStruct(toml347+"Num2JSON.toml", "CfgN2J347", pkgName, cfgSrc)

	strugen.GenNewCfg(cfgSrc)
}
