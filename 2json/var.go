package cvt2json

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cdutwhu/gotil/dispatcher"
	"github.com/cdutwhu/gotil/io"
	"github.com/cdutwhu/gotil/iter"
	"github.com/cdutwhu/gotil/misc"
	"github.com/cdutwhu/gotil/net"
	"github.com/cdutwhu/gotil/str"
	jt "github.com/cdutwhu/json-tool"
	"github.com/cdutwhu/n3-util/jkv"
	xt "github.com/cdutwhu/xml-tool"
	"github.com/digisan/logkit"
)

var (
	fPf              = fmt.Printf
	fPln             = fmt.Println
	fSp              = fmt.Sprint
	fSf              = fmt.Sprintf
	sHasPrefix       = strings.HasPrefix
	sHasSuffix       = strings.HasSuffix
	sReplace         = strings.Replace
	sReplaceAll      = strings.ReplaceAll
	sToLower         = strings.ToLower
	sTrim            = strings.Trim
	sCount           = strings.Count
	sSplit           = strings.Split
	sNewReader       = strings.NewReader
	sJoin            = strings.Join
	rxMustCompile    = regexp.MustCompile
	failOnErr        = logkit.FailOnErr
	enableLog2F      = logkit.Log2F
	failOnErrWhen    = logkit.FailOnErrWhen
	enableWarnDetail = logkit.WarnDetail
	warnOnErr        = logkit.WarnOnErr
	warner           = logkit.Warn
	localIP          = net.LocalIP
	splitRev         = str.SplitRev
	replByPosGrp     = str.ReplByPosGrp
	rmTailFromLastN  = str.RmTailFromLastN
	rmTailFromLast   = str.RmTailFromLast
	rmHeadToLast     = str.RmHeadToLast
	syncParallel     = dispatcher.SyncParallel
	iter2Slc         = iter.Iter2Slc
	mustWriteFile    = io.MustWriteFile
	trackTime        = misc.TrackTime
	xmlRoot          = xt.Root
	jsonRoot         = jt.Root
	fmtJSON          = jt.Fmt
	newJKV           = jkv.NewJKV
)

var (
	rxLB         = rxMustCompile(`\[[ \t\r\n]*\[`)
	rxRB         = rxMustCompile(`\][ \t\r\n]*\]`)
	rxOneEmpty   = rxMustCompile(`": \{\n([ ]+"-.+": .+,\n)*([ ]+"-.+": .+\n)[ ]+\}`)         // one empty object
	rxEmptyInArr = rxMustCompile(`[\[,]\n[ ]+\{\n([ ]+"-.+": .+,\n)*([ ]+"-.+": .+\n)[ ]+\}`) // empty object in array
	rxRawAttr    = rxMustCompile(`"-\w+":\s*`)                                                // raw style from XML attributes
)
