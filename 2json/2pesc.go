package cvt2json

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/digisan/gotk/slice/tso"
	jt "github.com/digisan/json-tool"
	"github.com/tidwall/gjson"
)

var (
	rRmPref = regexp.MustCompile(`"@\w+":`)
	mRepl   = map[string]string{
		`"#content"`: `"value"`,
	}
)

func Cvt2Pesc(goesJS, auditPath string) string {

	js := rRmPref.ReplaceAllStringFunc(goesJS, func(s string) string {
		return "\"" + s[2:]
	})
	for old, new := range mRepl {
		js = strings.ReplaceAll(js, old, new)
	}

	data := []byte(js)

	// audit
	// os.WriteFile(auditPath, data, os.ModePerm)

	///////////////////////////////////////////////////////////////////

	m := make(map[string]interface{})
	json.Unmarshal(data, &m)

	root, subs := "", []string{}
	if len(m) == 1 { // check this is NOT SIF Object
		for k, v := range m {
			root = k
			switch val := v.(type) {
			case map[string]interface{}:
				subs, _ = tso.Map2KVs(val, nil, nil)
			}
		}
	} else { // MUST be all SIF Object
		for object := range m {
			fmt.Printf("check [%s] must be SIF object\n", object)
		}
	}

	arrEles, objects := []string{}, []string{}

	for _, sub := range subs {
		path := fmt.Sprintf("%s.%s", root, sub) // sif.object
		// fmt.Println(path)

		if rst := gjson.Get(js, path); rst.IsArray() { // sif.array
			// fmt.Println("array SIF type:", sub, "-----", path)
			// fmt.Println(rst.String())

			for i := range rst.Array() {
				path := fmt.Sprintf("%s.%d", path, i)           // sif.array.0/1/2...
				if rst := gjson.Get(js, path); rst.IsObject() { // array element

					obj := fmt.Sprintf(`{"%s":%s}`, sub, rst.String()) // make json block from array element
					arrEles = append(arrEles, obj)
					// fmt.Println(obj)

					if !jt.IsValid([]byte(obj)) {
						panic("NOT VALID JSON")
					}

					// str, _ = sjson.Set(str, path, rst.String()) // 'sjson.Set' combines array, so DO NOT use it
					// fmt.Println(str)
				}
			}

		} else {
			// fmt.Println("other SIF type:", sub, "-----", path)

			rst := gjson.Get(js, path)
			if block := rst.String(); block[0] == '{' {
				obj := fmt.Sprintf(`{"%s":%s}`, sub, block) // make json block from array element
				if !jt.IsValid([]byte(obj)) {
					panic("NOT VALID JSON")
					// continue
				}
				objects = append(objects, obj)
			}
		}
	}

	// fmt.Println("all array objects count:", len(arrEles))

	objects = append(objects, arrEles...)
	pesc := fmt.Sprintf("[%s]", strings.Join(objects, ","))
	js = jt.FmtStr(pesc, "  ")

	if auditPath != "" {
		os.WriteFile(auditPath, []byte(js), os.ModePerm)
	}

	return js
}
