package cfg

import "github.com/cdutwhu/gonfig"

// Config : AUTO Created From /home/qmiao/Desktop/sif-xml2json/config/config.toml
type Config struct {
	Version interface{}
	Log string
	Service interface{}
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
		Timeout int
		URL string
		Subject string
	}
	Server struct {
		Protocol string
		IP interface{}
		Port interface{}
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
