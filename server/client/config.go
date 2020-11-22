package client

import "github.com/cdutwhu/gonfig"

// Config : AUTO Created From /home/qmiao/Desktop/sif-xml2json/server/client/config.toml
type Config struct {
	Service string
	Route struct {
		Help string
		Convert string
	}
	Server struct {
		Port int
		Protocol string
		IP string
	}
	Access struct {
		Timeout int
	}
}

// NewCfg :
func NewCfg(cfgStruName string, mReplExpr map[string]string, cfgPaths ...string) interface{} {
	var cfg interface{}
	switch cfgStruName {
	case "Config":
		cfg = &Config{}
	default:
		return nil
	}
	return gonfig.InitEnvVar(cfg, mReplExpr, cfgStruName, cfgPaths...)
}
