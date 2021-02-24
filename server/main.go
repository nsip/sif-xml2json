package main

import (
	"context"
	"io"
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
	jt "github.com/cdutwhu/json-tool"
	xt "github.com/cdutwhu/xml-tool"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nats-io/nats.go"
	sr "github.com/nsip/sif-spec-res"
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

func mkCfg4Docker(cfg interface{}) {
	forel := "./config_rel.toml"
	gonfig.Save(forel, cfg)
	outToml := "../config_d.toml"
	os.Remove(outToml)
	attrim.RmCfgAttrL1(forel, outToml, "Log", "Server", "Access")
}

var (
	gCfg *cfg.Config
)

func main() {
	// Running Executable, Single instance
	// if len(os.Args) == 1 {
	// 	one, err := single.New("sif-xml2json", single.WithLockPath("/tmp"))
	// 	failOnErr("%v", err)
	// 	failOnErr("%v", one.Lock())
	// 	defer func() {
	// 		failOnErr("%v", one.Unlock())
	// 		fPln("SIF-XML2JSON SERVER EXIT")
	// 	}()
	// }

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
	gCfg = pCfg.(*cfg.Config)

	// Trim a shorter config toml file for docker & client package
	if len(os.Args) > 2 && os.Args[2] == "trial" { // Args[1] == "--"
		mkCfg4Docker(gCfg)
		mkCfg4Clt(gCfg)
		return
	}

	ws := gCfg.WebService
	var IService interface{} = gCfg.Service // gCfg.Service can be "string" or "interface{}"
	service := IService.(string)

	// Set Jaeger Env for tracing
	os.Setenv("JAEGER_SERVICE_NAME", service)
	os.Setenv("JAEGER_SAMPLER_TYPE", "const")
	os.Setenv("JAEGER_SAMPLER_PARAM", "1")

	// Set LOGGLY
	setLoggly(false, gCfg.Loggly.Token, service)

	// Set Log Options
	syncBindLog(true)
	enableWarnDetail(false)
	if gCfg.Log != "" {
		enableLog2F(true, gCfg.Log)
		logGrp.Do(fSf("local log file @ [%s]", gCfg.Log))
	}

	logGrp.Do(fSf("[%s] Hosting on: [%v:%d], version [%v]", service, localIP(), ws.Port, gCfg.Version))

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
		// Cfg    = rflx.Env2Struct("Config", &cfg.Config{}).(*cfg.Config)
		port   = gCfg.WebService.Port
		fullIP = localIP() + fSf(":%d", port)
		route  = gCfg.Route
		mMtx   = initMutex(&gCfg.Route)
		vers   = sr.GetAllVer("v", "")
	)

	defer e.Start(fSf(":%d", port))
	logGrp.Do("Echo Service is Starting ...")

	// *************************************** List all API, FILE *************************************** //

	path := route.Help
	e.GET(path, func(c echo.Context) error {
		defer mMtx[path].Unlock()
		mMtx[path].Lock()

		return c.String(http.StatusOK,
			fSf("Converter Server Version: %s\n\n", gCfg.Version)+
				fSf("API:\n\n [POST] %-40s\n%s", fullIP+route.Convert,
					"\n Description: Upload SIF(XML), return SIF(JSON)\n"+
						"\n Parameters:"+
						"\n -- [sv]:   available SIF Ver: "+fSf("%v", vers)+
						"\n -- [nats]: send json to NATS?"+
						"\n -- [wrap]: is uploaded SIF file with single wrapped root?"))
	})

	// ------------------------------------------------------------------------------------ //

	// mRouteRes := map[string]string{
	// 	"/client-linux64": gCfg.File.ClientLinux64,
	// 	"/client-mac":     gCfg.File.ClientMac,
	// 	"/client-win64":   gCfg.File.ClientWin64,
	// 	"/client-config":  gCfg.File.ClientConfig,
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
		bytes, err := io.ReadAll(c.Request().Body)
		xstr, root, cont, lvl0 := "", "", "", ""
		xmlObjNames, xmlObjGrp, mObjContGrp := []string{}, []string{}, make(map[string][]string)

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
		// 	os.WriteFile("./debug.xml", []byte(xstr), 0666)
		// }
		/// DEBUG ///

		///
		// ** if wrapped, break and handle each SIF object ** //
		///
		root, lvl0, cont = xt.Lvl0(xstr)
		xmlObjNames, xmlObjGrp = []string{root}, []string{xstr}

		if wrapped {
			xmlObjNames, xmlObjGrp = xt.BreakCont(cont)
			jsonBuf, err := xj.Convert(sNewReader(lvl0))
			failOnErr("%v", err)
			lvl0json := jsonBuf.String()
			lvl0json = sReplaceAll(lvl0json, `""}`, `{`) // wrapper root without attributes
			lvl0json = sReplaceAll(lvl0json, `}}`, `,`)  // wrapper root with attributes
			RetSB.WriteString(lvl0json)
			RetSB.WriteString("\n")
		}
		///

		for i, xmlObj := range xmlObjGrp {
			obj := xmlObjNames[i]
			// logGrp.Do("cvt.XML2JSON")

			/// DEBUG ///
			// if sContains(obj, "Document") {
			// 	os.WriteFile("./debug.xml", []byte(xmlObj), 0666)
			// }
			/// DEBUG ///

			////////// ------------------- //////////

			// jsonObj, svApplied, err := cvt.XML2JSON(xmlObj, sv, false)
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
				return cvt.XML2JSON(xmlObj, sv, false)
			})

			jsonObj := results[0].Interface().(string)
			if !results[2].IsNil() {
				status = http.StatusInternalServerError
				RetSB.Reset()
				RetSB.WriteString(results[2].Interface().(error).Error())
				goto RET
			}
			logGrp.Do(obj + ": " + results[1].Interface().(string) + " applied")

			////////// ------------------- //////////

			if wrapped {
				_, jc := jt.SglEleBlkCont(jsonObj)
				mObjContGrp[obj] = append(mObjContGrp[obj], jc)
			} else {
				RetSB.WriteString(jsonObj)
				RetSB.WriteString("\n")
			}

			// Send a copy to NATS
			if msg {
				url, subj, timeout := gCfg.NATS.URL, gCfg.NATS.Subject, time.Duration(gCfg.NATS.Timeout)
				nc, err := nats.Connect(url)
				if err != nil {
					status = http.StatusInternalServerError
					RetSB.Reset()
					RetSB.WriteString(err.Error() + fSf(" @NATS Connect @Subject: [%s@%s]", url, subj))
					goto RET
				}
				msg, err := nc.Request(subj, []byte(jsonObj), timeout*time.Millisecond)
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
				Ret = sReplace(Ret, "\"-", "\"@", 1) // wrapper attribute '-' to '@'
				Ret = jt.Fmt(Ret, "  ")
			} else {
				Ret = RetSB.String()
			}
			logGrp.Do("--> Finish XML2JSON")
		}

		return c.String(status, sTrimRight(Ret, "\n")+"\n") // If already JSON String, so return String
	})
}
