package main

import (
	"flag"
	"io/ioutil"

	log "github.com/cihub/seelog"
	"github.com/getsentry/raven-go"
	"gopkg.in/yaml.v2"
)

var configfile = flag.String("c", "config.yaml", "The config file.")
var cfg struct {
	Sentry  string
	Server  string
	Logger  string
	DevMode bool
}
var Raven *raven.Client

func main() {
	parseCmdLine()

	initLogger()

	httpServer()
}

func parseCmdLine() {
	flag.Parse()
	if conf, err := ioutil.ReadFile(*configfile); err != nil {
		panic(err)
	} else {
		if err = yaml.Unmarshal(conf, &cfg); err != nil {
			panic(err)
		}
	}
}

func initLogger() {
	var err error

	// 初始化raven
	if cfg.Sentry != "" {
		Raven, err = raven.NewClient(cfg.Sentry, nil)
		if err != nil {
			log.Error("Init Sentry Error:", err)
		}
	}

	// 初始化logger
	if cfg.Logger != "" {
		// 自定义一个seelog raven receiver
		receiver := &RavenReciver{Client: Raven}
		parseParams := &log.CfgParseParams{
			CustomReceiverProducers: map[string]log.CustomReceiverProducer{
				"sentry": func(log.CustomReceiverInitArgs) (log.CustomReceiver, error) {
					return receiver, nil
				},
			},
		}
		if logger, err := log.LoggerFromParamConfigAsFile(cfg.Logger, parseParams); err == nil {
			log.ReplaceLogger(logger)
		} else {
			log.Error("Parse Logger Error: ", err)
		}

	}

}
