package cfg

import "github.com/cdutwhu/gonfig"

// Config : AUTO Created From "sif-xml2json/config/config.toml"
type Config struct {
	Log string
	Service interface{}
	Version interface{}
	SIF struct {
		DefaultVer string
	}
	Loggly struct {
		Token string
	}
	WebService struct {
		Port int
	}
	Route struct {
		Convert string
		Help string
	}
	NATS struct {
		URL string
		Subject string
		Timeout int
	}
	Server struct {
		Port interface{}
		Protocol string
		IP interface{}
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
