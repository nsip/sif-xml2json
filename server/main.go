package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"time"

	xj "github.com/basgys/goxml2json"
	"github.com/cdutwhu/gonfig"
	"github.com/cdutwhu/gonfig/attrim"
	"github.com/cdutwhu/gonfig/strugen"
	"github.com/cdutwhu/gotil/misc"
	"github.com/cdutwhu/gotil/rflx"
	jt "github.com/cdutwhu/json-tool"
	xt "github.com/cdutwhu/xml-tool"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nats-io/nats.go"
	cvt "github.com/nsip/sif-xml2json/2json"
	cfg "github.com/nsip/sif-xml2json/config/cfg"
	errs "github.com/nsip/sif-xml2json/err-const"
)

func mkCfg4Clt(cfg interface{}) {
	forel := "./config_rel.toml"
	gonfig.Save(forel, cfg)
	outToml := "./client/config.toml"
	outSrc := "./client/config.go"
	os.Remove(outToml)
	os.Remove(outSrc)
	attrim.SelCfgAttrL1(forel, outToml, "Service", "Route", "Server", "Access")
	strugen.GenStruct(outToml, "Config", "client", outSrc)
	strugen.GenNewCfg(outSrc)
}

func main() {
	// Load global config.toml file from config/
	gonfig.SetDftCfgVal("sif-xml2json", "0.0.0")
	pCfg := cfg.NewCfg(
		"Config",
		map[string]string{
			"[s]":    "Service",
			"[v]":    "Version",
			"[port]": "WebService.Port",
		},
		"./config/config.toml",
		"../config/config.toml",
	)
	failOnErrWhen(pCfg == nil, "%v: Config Init Error", errs.CFG_INIT_ERR)
	Cfg := pCfg.(*cfg.Config)

	// Trim a shorter config toml file for client package
	if len(os.Args) > 2 && os.Args[2] == "trial" {
		mkCfg4Clt(Cfg)
		return
	}

	ws := Cfg.WebService
	var IService interface{} = Cfg.Service // Cfg.Service can be "string", can be "interface{}"
	service := IService.(string)

	// Set Jaeger Env for tracing
	os.Setenv("JAEGER_SERVICE_NAME", service)
	os.Setenv("JAEGER_SAMPLER_TYPE", "const")
	os.Setenv("JAEGER_SAMPLER_PARAM", "1")

	// Set LOGGLY
	setLoggly(true, Cfg.Loggly.Token, service)

	// Set Log Options
	syncBindLog(true)
	enableWarnDetail(false)
	enableLog2F(true, Cfg.Log)

	logGrp.Do(fSf("local log file @ [%s]", Cfg.Log))
	logGrp.Do(fSf("[%s] Hosting on: [%v:%d], version [%v]", service, localIP(), ws.Port, Cfg.Version))

	// Start Service
	done := make(chan string)
	c := make(chan os.Signal)
	signal.Notify(c, os.Kill, os.Interrupt)
	go HostHTTPAsync(c, done)
	logGrp.Do(<-done)
}

func shutdownAsync(e *echo.Echo, sig <-chan os.Signal, done chan<- string) {
	<-sig
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	failOnErr("%v", e.Shutdown(ctx))
	time.Sleep(20 * time.Millisecond)
	done <- "Shutdown Successfully"
}

// HostHTTPAsync : Host a HTTP Server for XML to JSON
func HostHTTPAsync(sig <-chan os.Signal, done chan<- string) {
	defer logGrp.Do("HostHTTPAsync Exit")

	e := echo.New()
	defer e.Close()

	// waiting for shutdown
	go shutdownAsync(e, sig, done)

	// Add Jaeger Tracer into Middleware
	c := jaegertracing.New(e, nil)
	defer c.Close()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("2G"))
	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.POST},
		AllowCredentials: true,
	}))

	e.Logger.SetOutput(os.Stdout)
	e.Logger.Infof(" ------------------------ e.Logger.Infof ------------------------ ")

	var (
		Cfg    = rflx.Env2Struct("Config", &cfg.Config{}).(*cfg.Config)
		port   = Cfg.WebService.Port
		fullIP = localIP() + fSf(":%d", port)
		route  = Cfg.Route
		mMtx   = initMutex(&Cfg.Route)
	)

	defer e.Start(fSf(":%d", port))
	logGrp.Do("Echo Service is Starting ...")

	// *************************************** List all API, FILE *************************************** //

	path := route.Help
	e.GET(path, func(c echo.Context) error {
		defer mMtx[path].Unlock()
		mMtx[path].Lock()

		return c.String(http.StatusOK,
			// fSf("wget %-55s-> %s\n", fullIP+"/client-linux64", "Get Client(Linux64)")+
			// 	fSf("wget %-55s-> %s\n", fullIP+"/client-mac", "Get Client(Mac)")+
			// 	fSf("wget %-55s-> %s\n", fullIP+"/client-win64", "Get Client(Windows64)")+
			// 	fSf("wget -O config.toml %-40s-> %s\n", fullIP+"/client-config", "Get Client Config")+
			// 	fSf("\n")+
			fSf("[POST] %-40s\n%s", fullIP+route.Convert,
				"--- Upload SIF(XML), return SIF(JSON).\n"+
					"------ [sv]:   SIF Ver\n"+
					"------ [nats]: send json to NATS\n"+
					"------ [wrap]: if uploaded SIF is single root wrapped file"))
	})

	// ------------------------------------------------------------------------------------ //

	// mRouteRes := map[string]string{
	// 	"/client-linux64": Cfg.File.ClientLinux64,
	// 	"/client-mac":     Cfg.File.ClientMac,
	// 	"/client-win64":   Cfg.File.ClientWin64,
	// 	"/client-config":  Cfg.File.ClientConfig,
	// }

	// routeFun := func(rt, res string) func(c echo.Context) error {
	// 	return func(c echo.Context) (err error) {
	// 		if _, err = os.Stat(res); err == nil {
	// 			fPln(rt, res)
	// 			return c.File(res)
	// 		}
	// 		return warnOnErr("FILE_NOT_FOUND: [%s]  get [%s]", rt, res)
	// 	}
	// }

	// for rt, res := range mRouteRes {
	// 	e.GET(rt, routeFun(rt, res))
	// }

	// -------------------------------------------------------------------------------- //
	// -------------------------------------------------------------------------------- //

	path = route.Convert
	e.POST(path, func(c echo.Context) error {
		defer misc.TrackTime(time.Now())
		defer mMtx[path].Unlock()
		mMtx[path].Lock()

		var (
			status  = http.StatusOK
			Ret     = ""
			RetSB   strings.Builder
			results []reflect.Value // for 'jaegertracing.TraceFunction'
		)

		logGrp.Do("Parsing Params")
		pvalues, sv, msg, wrapped := c.QueryParams(), "", false, false
		if ok, v := url1Value(pvalues, 0, "sv"); ok {
			sv = v
		}
		if ok, n := url1Value(pvalues, 0, "nats"); ok && n != "false" {
			msg = true
		}
		if ok, w := url1Value(pvalues, 0, "wrap"); ok && w != "false" {
			wrapped = true
		}

		logGrp.Do("Reading Request Body")
		bytes, err := ioutil.ReadAll(c.Request().Body)
		xstr, root, cont, lvl0 := "", "", "", ""
		sifObjNames, sifObjGrp, mObjContGrp := []string{}, []string{}, make(map[string][]string)

		if err != nil {
			status = http.StatusInternalServerError
			RetSB.Reset()
			RetSB.WriteString(err.Error() + " @Request Body")
			goto RET
		}
		if xstr = string(bytes); len(xstr) == 0 {
			status = http.StatusBadRequest
			RetSB.Reset()
			RetSB.WriteString("HTTP_REQBODY_EMPTY @Request Body")
			goto RET
		}
		if !isXML(xstr) {
			status = http.StatusBadRequest
			RetSB.Reset()
			RetSB.WriteString("PARAM_INVALID_XML @Request Body")
			goto RET
		}

		/// DEBUG ///
		// if sContains(xstr, "A5A575C7-8917-5101-B8E7-F08ED123A823") {
		// 	ioutil.WriteFile("./debug.xml", []byte(xstr), 0666)
		// }
		/// DEBUG ///

		///
		// ** if wrapped, break and handle each SIF object ** //
		///
		root, lvl0, cont = xt.Lvl0(xstr)
		sifObjNames, sifObjGrp = []string{root}, []string{xstr}

		if wrapped {
			sifObjNames, sifObjGrp = xt.BreakCont(cont)
			jsonBuf, err := xj.Convert(sNewReader(lvl0))
			failOnErr("%v", err)
			lvl0json := jsonBuf.String()
			lvl0json = sReplaceAll(lvl0json, `""}`, `{`) // wrapper root without attributes
			lvl0json = sReplaceAll(lvl0json, `}}`, `,`)  // wrapper root with attributes
			RetSB.WriteString(lvl0json)
			RetSB.WriteString("\n")
		}
		///

		for i, objsif := range sifObjGrp {
			obj := sifObjNames[i]
			// logGrp.Do("cvt.XML2JSON")

			/// DEBUG ///
			// if sContains(obj, "Document") {
			// 	ioutil.WriteFile("./debug.xml", []byte(objsif), 0666)
			// }
			/// DEBUG ///

			////////// ------------------- //////////

			// objson, svApplied, err := cvt.XML2JSON(objsif, sv, false)
			// if err != nil {
			// 	status = http.StatusInternalServerError
			// 	RetSB.Reset()
			// 	RetSB.WriteString(err.Error())
			// 	goto RET
			// }
			// logGrp.Do(obj + ": " + svApplied + " applied")

			////////// ------------------- //////////

			// Trace [cvt.XML2JSON], uses (variadic parameter), must wrap it to [jaegertracing.TraceFunction]
			results = jaegertracing.TraceFunction(c, func() (string, string, error) {
				return cvt.XML2JSON(objsif, sv, false)
			})
			objson := results[0].Interface().(string)
			if !results[2].IsNil() {
				status = http.StatusInternalServerError
				RetSB.Reset()
				RetSB.WriteString(results[2].Interface().(error).Error())
				goto RET
			}
			logGrp.Do(obj + ": " + results[1].Interface().(string) + " applied")

			////////// ------------------- //////////

			if wrapped {
				_, jc := jt.SglEleBlkCont(objson)
				mObjContGrp[obj] = append(mObjContGrp[obj], jc)
			} else {
				RetSB.WriteString(objson)
				RetSB.WriteString("\n")
			}

			// Send a copy to NATS
			if msg {
				url, subj, timeout := Cfg.NATS.URL, Cfg.NATS.Subject, time.Duration(Cfg.NATS.Timeout)
				nc, err := nats.Connect(url)
				if err != nil {
					status = http.StatusInternalServerError
					RetSB.Reset()
					RetSB.WriteString(err.Error() + fSf(" @NATS Connect @Subject: [%s@%s]", url, subj))
					goto RET
				}
				msg, err := nc.Request(subj, []byte(objson), timeout*time.Millisecond)
				if err != nil {
					status = http.StatusInternalServerError
					RetSB.Reset()
					RetSB.WriteString(err.Error() + fSf(" @NATS Request @Subject: [%s@%s]", url, subj))
					goto RET
				}
				logGrp.Do(string(msg.Data))
			}
		}

	RET:
		if status != http.StatusOK {
			Ret = RetSB.String()
			warnGrp.Do(Ret + " --> Failed")
		} else {
			if wrapped {
				for obj, conts := range mObjContGrp {
					if len(conts) == 1 {
						RetSB.WriteString(fSf(`"%s": %s,`, obj, conts[0]))
					} else {
						RetSB.WriteString(fSf(`"%s": [%s],`, obj, sJoin(conts, ",")))
					}
				}
				Ret = RetSB.String()
				Ret = sTrimRight(Ret, ",") + "}}"
				Ret = jt.Fmt(Ret, "  ")
			} else {
				Ret = RetSB.String()
			}
			logGrp.Do("--> Finish XML2JSON")
		}

		return c.String(status, sTrimRight(Ret, "\n")+"\n") // If already JSON String, so return String
	})
}
