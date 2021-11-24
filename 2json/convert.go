package cvt2json

import (
	"fmt"

	xj "github.com/basgys/goxml2json"
	"github.com/digisan/gotk/io"
	sif342 "github.com/nsip/sif-spec-res/3.4.2"
	sif343 "github.com/nsip/sif-spec-res/3.4.3"
	sif344 "github.com/nsip/sif-spec-res/3.4.4"
	sif345 "github.com/nsip/sif-spec-res/3.4.5"
	sif346 "github.com/nsip/sif-spec-res/3.4.6"
	sif347 "github.com/nsip/sif-spec-res/3.4.7"
	sif348 "github.com/nsip/sif-spec-res/3.4.8"
	sif349 "github.com/nsip/sif-spec-res/3.4.9"
	cfg "github.com/nsip/sif-xml2json/config/cfg"
)

func AllSIFObject(ver string) (objs []string, err error) {
	txt, err := BytesOfTXT(ver)
	if err != nil {
		return nil, err
	}
	const pfx = "OBJECT: "
	p := len(pfx)
	objstr, err := io.StrLineScan(txt, func(line string) (bool, string) {
		if sHasPrefix(line, pfx) {
			return true, line[p:]
		}
		return false, ""
	}, "")
	return sSplit(objstr, "\n"), err
}

func BytesOfTXT(ver string) (ret string, err error) {
	var mBytes map[string][]byte
	switch ver {
	case "3.4.2":
		mBytes = sif342.TXT
	case "3.4.3":
		mBytes = sif343.TXT
	case "3.4.4":
		mBytes = sif344.TXT
	case "3.4.5":
		mBytes = sif345.TXT
	case "3.4.6":
		mBytes = sif346.TXT
	case "3.4.7":
		mBytes = sif347.TXT
	case "3.4.8":
		mBytes = sif348.TXT
	case "3.4.9":
		mBytes = sif349.TXT
	default:
		err = fmt.Errorf("error: No SIF Spec @ Version [%s]", ver)
		warner("%v", err)
		return
	}
	return string(mBytes[sReplaceAll(ver, ".", "")]), err
}

func BytesOfJSON(ver, ruleType, object string, indices ...int) (ret []string, err error) {

	var mBytes map[string][]byte
	switch ver {
	case "3.4.2":
		switch sToLower(ruleType) {
		case "bool", "boolean":
			mBytes = sif342.JSON_BOOL
		case "list":
			mBytes = sif342.JSON_LIST
		case "num", "number", "numeric":
			mBytes = sif342.JSON_NUM
		}
	case "3.4.3":
		switch sToLower(ruleType) {
		case "bool", "boolean":
			mBytes = sif343.JSON_BOOL
		case "list":
			mBytes = sif343.JSON_LIST
		case "num", "number", "numeric":
			mBytes = sif343.JSON_NUM
		}
	case "3.4.4":
		switch sToLower(ruleType) {
		case "bool", "boolean":
			mBytes = sif344.JSON_BOOL
		case "list":
			mBytes = sif344.JSON_LIST
		case "num", "number", "numeric":
			mBytes = sif344.JSON_NUM
		}
	case "3.4.5":
		switch sToLower(ruleType) {
		case "bool", "boolean":
			mBytes = sif345.JSON_BOOL
		case "list":
			mBytes = sif345.JSON_LIST
		case "num", "number", "numeric":
			mBytes = sif345.JSON_NUM
		}
	case "3.4.6":
		switch sToLower(ruleType) {
		case "bool", "boolean":
			mBytes = sif346.JSON_BOOL
		case "list":
			mBytes = sif346.JSON_LIST
		case "num", "number", "numeric":
			mBytes = sif346.JSON_NUM
		}
	case "3.4.7":
		switch sToLower(ruleType) {
		case "bool", "boolean":
			mBytes = sif347.JSON_BOOL
		case "list":
			mBytes = sif347.JSON_LIST
		case "num", "number", "numeric":
			mBytes = sif347.JSON_NUM
		}
	case "3.4.8":
		switch sToLower(ruleType) {
		case "bool", "boolean":
			mBytes = sif348.JSON_BOOL
		case "list":
			mBytes = sif348.JSON_LIST
		case "num", "number", "numeric":
			mBytes = sif348.JSON_NUM
		}
	case "3.4.9":
		switch sToLower(ruleType) {
		case "bool", "boolean":
			mBytes = sif349.JSON_BOOL
		case "list":
			mBytes = sif349.JSON_LIST
		case "num", "number", "numeric":
			mBytes = sif349.JSON_NUM
		}
	default:
		err = fmt.Errorf("Error: No SIF Spec @ Version [%s]", ver)
		warner("%v", err)
		return
	}

	for _, idx := range indices {
		key := fSf("%s_%d", object, idx)
		if bytes, ok := mBytes[key]; ok {
			ret = append(ret, string(bytes))
		}
	}
	return
}

// enforceCfg : LIST config must be from low Level to high level
func enforceCfg(json string, lsJSONCfg ...string) string {
	for _, jsoncfg := range lsJSONCfg {
		// make sure [jsoncfg] is formatted; Otherwise, do Fmt firstly
		// jsoncfg = fmtJSON(jsoncfg, "  ")

		json, _ = newJKV(json, "", false).Unfold(0, newJKV(jsoncfg, "", false))
		// make sure there is no double "[" OR "]"
		bytes := rxRB.ReplaceAll(rxLB.ReplaceAll([]byte(json), []byte("[")), []byte("]"))
		json = fmtJSON(string(bytes), "  ")
	}
	return json
}

// XML2JSON : if [sifver] is "", DefaultSIFVer applies
func XML2JSON(xml, sifver string, enforced bool, subobj ...string) (string, string, error) {
	cfgAll := cfg.NewCfg("Config", nil, "./config/config.toml", "../config/config.toml").(*cfg.Config)

	jsonBuf, err := xj.Convert(
		sNewReader(xml),
		// xj.WithTypeConverter(xj.Float, xj.Int, xj.Bool, xj.Null),
		// xj.WithAttrPrefix("-"),
		// xj.WithContentPrefix("#"),
	)
	failOnErr("That's embarrassing... %v", err)

	// json := jsonBuf.String()
	// return // --------------------------- test 3rd party lib --------------------------- //

	json := fmtJSON(jsonBuf.String(), "  ")

	// Deal with 'LF', 'TB', Part1 -------------------------------------------------------- //
	mRepl1 := map[string]string{"\n": "#LF#", "\t": "#TB#"}
	for k, v := range mRepl1 {
		posGrp, values := [][]int{}, []string{}
		re := rxMustCompile(fSf(`": "[^"]*[%s]+[^"]*"[,\n]{1}`, k))
		for _, pos := range re.FindAllStringIndex(json, -1) {
			start, end := pos[0]+4, pos[1]-2
			posGrp = append(posGrp, []int{start, end})
			values = append(values, sReplaceAll(json[start:end], k, v))
		}
		json = replByPosGrp(json, posGrp, values)
	}

	// Attributes Modification according to Config ---------------------------------------- //
	obj := xmlRoot(xml)              // infer object from xml root, use this object to find config json by default
	if enforced && len(subobj) > 0 { // if object is provided, ignore default, use 1st given object to search
		obj = subobj[0]
	}

	ver := cfgAll.SIF.DefaultVer
	if sifver != "" {
		ver = sifver
	}

	if rt, err := BytesOfJSON(ver, "list", obj, iter2Slc(10)...); err == nil {
		json = enforceCfg(json, rt...)
	} else {
		return "", "", err
	}
	if rt, err := BytesOfJSON(ver, "num", obj, iter2Slc(2)...); err == nil {
		json = enforceCfg(json, rt...)
	} else {
		return "", "", err
	}
	if rt, err := BytesOfJSON(ver, "bool", obj, iter2Slc(2)...); err == nil {
		json = enforceCfg(json, rt...)
	} else {
		return "", "", err
	}

	// Deal with 'LF', 'TB'  Part2 -------------------------------------------------------------
	mRepl2 := map[string]string{"#LF#": "\\n", "#TB#": "\\t"}
	for k, v := range mRepl2 {
		json = sReplaceAll(json, k, v)
	}

	// XML empty element(empty text) with Attributes -------------------------------------------
	emptyPosPair := [][]int{}
	for _, pos := range rxOneEmpty.FindAllStringIndex(json, -1) {
		emptyPosPair = append(emptyPosPair, []int{pos[0] + 6, pos[0] + 6})
	}
	for _, pos := range rxEmptyInArr.FindAllStringIndex(json, -1) {
		remain, offset := json[pos[0]:], 0
		for i, c := range remain {
			if c == '{' {
				offset = i + 1
				break
			}
		}
		emptyPosPair = append(emptyPosPair, []int{pos[0] + offset, pos[0] + offset})
	}
	const mark = "#content" // "value"
	json = replByPosGrp(json, emptyPosPair, []string{fSf("\"%s\": \"\",\n", mark)})

	// "-Attribute" => "@Attribute"
	json = rxRawAttr.ReplaceAllStringFunc(json, func(m string) string {
		return sReplace(m, `"-`, `"@`, 1)
	})

	json = fmtJSON(json, "  ")

	///////////////////////////////////
	// TO PESC

	


	///////////////////////////////////

	return json, ver, nil
}
