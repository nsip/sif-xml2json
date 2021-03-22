package cvt2json

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	errs "github.com/nsip/sif-xml2json/err-const"
)

func TestXMLRoot(t *testing.T) {
	bytes, err := os.ReadFile("../data/examples/3.4.8/Activity_0.xml")
	failOnErr("%v", err)
	fPln(xmlRoot(string(bytes)))
}

func x2j(dim int, tid int, done chan int, params ...interface{}) {
	defer func() { done <- tid }()

	files := params[0].([]os.FileInfo)
	dir := params[1].(string)
	ver := params[2].(string)

	for i := tid; i < len(files); i += dim {
		obj := rmTailFromLast(files[i].Name(), ".")
		fPln("start:", obj)
		// if exist(obj, "LearningStandardDocument", "StudentAttendanceTimeList") {
		// 	continue
		// }
		bytes, err := os.ReadFile(filepath.Join(dir, files[i].Name()))
		failOnErr("%v", err)
		json, sv, err := XML2JSON(string(bytes), ver, false)
		fPln("end:", obj, sv, err)
		if json != "" {
			mustWriteFile(fSf("../data/output/%s/%s.json", sv, obj), []byte(json))
		}
	}
}

func TestXML2JSON(t *testing.T) {
	defer trackTime(time.Now())
	// enableLog2F(true, "./error.log")
	// defer enableLog2F(false, "")
	// defer enableWarnDetail(true)
	enableWarnDetail(false)

	ver := "3.4.8"
	dir := `../data/examples/` + ver
	files, err := os.ReadDir(dir)
	failOnErr("%v", err)
	failOnErrWhen(len(files) == 0, "%v", errs.FILE_NOT_FOUND)
	syncParallel(1, x2j, files, dir, ver) // multi threads may be error.
	fPln("OK")
}
