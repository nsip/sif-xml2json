package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/io"
	"github.com/cdutwhu/gotil/judge"
	"github.com/cdutwhu/gotil/net"
	"github.com/cdutwhu/gotil/rflx"
	"github.com/cdutwhu/n3-util/n3log"
	"github.com/cdutwhu/n3-util/rest"
)

var (
	fPln             = fmt.Println
	fSf              = fmt.Sprintf
	fPf              = fmt.Printf
	sReplaceAll      = strings.ReplaceAll
	sTrimRight       = strings.TrimRight
	sTrimLeft        = strings.TrimLeft
	sTrim            = strings.Trim
	sIndex           = strings.Index
	sLastIndex       = strings.LastIndex
	sJoin            = strings.Join
	sNewReader       = strings.NewReader
	sHasSuffix       = strings.HasSuffix
	sHasPrefix       = strings.HasPrefix
	sContains        = strings.Contains
	rxMustCompile    = regexp.MustCompile
	failOnErr        = fn.FailOnErr
	failOnErrWhen    = fn.FailOnErrWhen
	enableLog2F      = fn.EnableLog2F
	enableWarnDetail = fn.EnableWarnDetail
	logWhen          = fn.LoggerWhen
	logger           = fn.Logger
	warnOnErr        = fn.WarnOnErr
	warner           = fn.Warner
	localIP          = net.LocalIP
	env2Struct       = rflx.Env2Struct
	struct2Map       = rflx.Struct2Map
	tryInvoke        = rflx.TryInvoke
	loggly           = n3log.Loggly
	logBind          = n3log.Bind
	setLoggly        = n3log.SetLoggly
	syncBindLog      = n3log.SyncBindLog
	isXML            = judge.IsXML
	isJSON           = judge.IsJSON
	exist            = judge.Exist
	mustWriteFile    = io.MustWriteFile
	url1Value        = rest.URL1Value
)

var (
	logGrp  = logBind(logger) // logBind(logger, loggly("info"))
	warnGrp = logBind(warner) // logBind(warner, loggly("warn"))
)

func initMutex(route interface{}) map[string]*sync.Mutex {
	mMtx := make(map[string]*sync.Mutex)
	for _, v := range struct2Map(route) {
		mMtx[v.(string)] = &sync.Mutex{}
	}
	return mMtx
}
