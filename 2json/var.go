package cvt2json

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/dispatcher"
	"github.com/cdutwhu/gotil/io"
	"github.com/cdutwhu/gotil/iter"
	"github.com/cdutwhu/gotil/judge"
	"github.com/cdutwhu/gotil/misc"
	"github.com/cdutwhu/gotil/net"
	"github.com/cdutwhu/gotil/str"
	jt "github.com/cdutwhu/json-tool"
	"github.com/cdutwhu/n3-util/jkv"
	xt "github.com/cdutwhu/xml-tool"
)

var (
	fPf              = fmt.Printf
	fPln             = fmt.Println
	fSp              = fmt.Sprint
	fSf              = fmt.Sprintf
	sHasPrefix       = strings.HasPrefix
	sHasSuffix       = strings.HasSuffix
	sReplaceAll      = strings.ReplaceAll
	sToLower         = strings.ToLower
	sTrim            = strings.Trim
	sCount           = strings.Count
	sSplit           = strings.Split
	sNewReader       = strings.NewReader
	sJoin            = strings.Join
	rxMustCompile    = regexp.MustCompile
	failOnErr        = fn.FailOnErr
	enableLog2F      = fn.EnableLog2F
	failOnErrWhen    = fn.FailOnErrWhen
	enableWarnDetail = fn.EnableWarnDetail
	warnOnErr        = fn.WarnOnErr
	warner           = fn.Warner
	localIP          = net.LocalIP
	splitRev         = str.SplitRev
	replByPosGrp     = str.ReplByPosGrp
	rmTailFromLastN  = str.RmTailFromLastN
	rmTailFromLast   = str.RmTailFromLast
	rmHeadToLast     = str.RmHeadToLast
	syncParallel     = dispatcher.SyncParallel
	iter2Slc         = iter.Iter2Slc
	mustWriteFile    = io.MustWriteFile
	exist            = judge.Exist
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
)
