package client

import "github.com/cdutwhu/gonfig"

// Config : AUTO Created From /sif-xml2json/server/client/config.toml
type Config struct {
	Service string
	Route struct {
		Help string
		Convert string
	}
	Server struct {
		Protocol string
		IP string
		Port int
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
