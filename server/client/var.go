package client

import (
	"fmt"
	"strings"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/io"
	"github.com/cdutwhu/gotil/judge"
	"github.com/cdutwhu/gotil/rflx"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

var (
	fPt           = fmt.Print
	fPf           = fmt.Printf
	fPln          = fmt.Println
	fSf           = fmt.Sprintf
	sJoin         = strings.Join
	sReplace      = strings.Replace
	sReplaceAll   = strings.ReplaceAll
	sTrimRight    = strings.TrimRight
	struct2Map    = rflx.Struct2Map
	mapKeys       = rflx.MapKeys
	env2Struct    = rflx.Env2Struct
	struct2Env    = rflx.Struct2Env
	enableLog2F   = fn.EnableLog2F
	failOnErrWhen = fn.FailOnErrWhen
	failOnErr     = fn.FailOnErr
	logWhen       = fn.LoggerWhen
	warnOnErr     = fn.WarnOnErr
	warnOnErrWhen = fn.WarnOnErrWhen
	isXML         = judge.IsXML
	isJSON        = judge.IsJSON
	mustWriteFile = io.MustWriteFile
)

// Args is arguments for "Route"
type Args struct {
	Data   []byte
	Ver    string
	ToNATS bool
	Wrap   bool
}

func initMapFnURL(protocol, ip string, port int, route interface{}) (map[string]string, []string) {
	mFnURL := make(map[string]string)
	for k, v := range struct2Map(route) {
		mFnURL[k] = fSf("%s://%s:%d%s", protocol, ip, port, v)
	}
	return mFnURL, mapKeys(mFnURL).([]string)
}

func initTracer(serviceName string) opentracing.Tracer {
	cfg, err := config.FromEnv()
	failOnErr("%v: ", err)
	cfg.ServiceName = serviceName
	cfg.Sampler.Type = "const"
	cfg.Sampler.Param = 1

	tracer, _, err := cfg.NewTracer()
	failOnErr("%v: ", err)
	return tracer
}
