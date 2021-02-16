package client

import "github.com/cdutwhu/gonfig"

// Config : AUTO Created From /sif-xml2json/server/client/config.toml
type Config struct {
	Service string
	Route struct {
		Convert string
		Help string
	}
	Server struct {
		IP string
		Port int
		Protocol string
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
