package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/cdutwhu/gotil/judge"
	"github.com/cdutwhu/gotil/net"
	"github.com/cdutwhu/gotil/rflx"
	"github.com/cdutwhu/n3-util/n3log"
	"github.com/cdutwhu/n3-util/rest"
	"github.com/digisan/logkit"
)

var (
	fSf              = fmt.Sprintf
	sReplace         = strings.Replace
	sReplaceAll      = strings.ReplaceAll
	sTrimRight       = strings.TrimRight
	sJoin            = strings.Join
	sNewReader       = strings.NewReader
	failOnErr        = logkit.FailOnErr
	failOnErrWhen    = logkit.FailOnErrWhen
	enableLog2F      = logkit.Log2F
	enableWarnDetail = logkit.WarnDetail
	logger           = logkit.Log
	warner           = logkit.Warn
	localIP          = net.LocalIP
	struct2Map       = rflx.Struct2Map
	logBind          = n3log.Bind
	setLoggly        = n3log.SetLoggly
	syncBindLog      = n3log.SyncBindLog
	isXML            = judge.IsXML
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
