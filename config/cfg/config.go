package cfg

import "github.com/cdutwhu/gonfig"

// Config : AUTO Created From /home/qmiao/Desktop/4rel/sif-xml2json/config/config.toml
type Config struct {
	Service interface{}
	Version interface{}
	Log string
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
		Subject string
		Timeout int
		URL string
	}
	Server struct {
		IP interface{}
		Port interface{}
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
