package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	errs "github.com/nsip/sif-xml2json/err-const"
	"github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"
)

// DoWithTrace :
func DoWithTrace(ctx context.Context, config, fn string, args *Args) (string, error) {
	pCfg := NewCfg("Config", nil)
	failOnErrWhen(pCfg == nil, "%v", errs.CFG_INIT_ERR)
	Cfg := pCfg.(*Config)

	service := Cfg.Service
	if span := opentracing.SpanFromContext(ctx); span != nil {
		tracer := initTracer(service)
		span := tracer.StartSpan(fn, opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, service)
		if args != nil {
			span.SetTag(fn, *args)
		}
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return DO(config, fn, args)
}

// DO : fn ["Help", "Convert"]
func DO(config, fn string, args *Args) (string, error) {
	pCfg := NewCfg("Config", nil)
	failOnErrWhen(pCfg == nil, "%v", errs.CFG_INIT_ERR)
	Cfg := pCfg.(*Config)

	server := Cfg.Server
	protocol, ip, port := server.Protocol, server.IP, server.Port
	timeout := Cfg.Access.Timeout

	mFnURL, fields := initMapFnURL(protocol, ip, port, &Cfg.Route)
	url, ok := mFnURL[fn]
	if !ok {
		warnOnErrWhen(!ok, "%v: Need %v", errs.PARAM_NOT_SUPPORTED, fields)
		return "", errs.PARAM_NOT_SUPPORTED
	}

	chStr, chErr := make(chan string), make(chan error)
	go func() {
		rest(fn, url, args, chStr, chErr)
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		warnOnErr("%v: Didn't get response in %d(s)", errs.NET_TIMEOUT, timeout)
		return "", errs.NET_TIMEOUT
	case str := <-chStr:
		err := <-chErr
		if err == errs.NO_ERROR {
			return str, nil
		}
		return str, err
	}
}

// rest :
func rest(fn, url string, args *Args, chStr chan string, chErr chan error) {
	pVer, pNats, pWrap := "", "", ""
	if args != nil {
		if args.Ver != "" {
			pVer = fSf("sv=%s", args.Ver)
		}
		if args.ToNATS {
			pNats = fSf("nats")
		}
		if args.Wrap {
			pWrap = fSf("wrap")
		}
	}

	url = fSf("%s?%s&%s&%s", url, pVer, pNats, pWrap)
	for i := 0; i < 16; i++ {
		url = sReplaceAll(url, "&&", "&") // remove empty params
	}
	url = sReplace(url, "?&", "?", 1)
	url = sTrimRight(url, "?&")

	logWhen(true, "accessing ... %s", url)

	var (
		Resp *http.Response
		Err  error
		Ret  []byte
	)

	switch fn {
	case "Help":
		if Resp, Err = http.Get(url); Err != nil {
			goto ERR_RET
		}
	case "Convert":
		if args == nil {
			Err = errs.PARAM_INVALID
			goto ERR_RET
		}
		if !isXML(string(args.Data)) {
			Err = errs.PARAM_INVALID_XML
			goto ERR_RET
		}
		if Resp, Err = http.Post(url, "application/json", bytes.NewBuffer(args.Data)); Err != nil {
			goto ERR_RET
		}
	default:
		panic("Shouldn't be here")
	}

	if Resp == nil {
		Err = errs.NET_NO_RESPONSE
		goto ERR_RET
	}
	defer Resp.Body.Close()

	if Ret, Err = io.ReadAll(Resp.Body); Err != nil {
		goto ERR_RET
	}

ERR_RET:
	if Err != nil {
		chStr <- ""
		chErr <- Err
		return
	}

	chStr <- string(Ret)
	chErr <- errs.NO_ERROR
}
