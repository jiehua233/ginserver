package main

import (
	"flag"
	"io/ioutil"

	"github.com/getsentry/raven-go"
	"gopkg.in/yaml.v2"
)

var configfile = flag.String("c", "config.yaml", "The config file.")
var cfg struct {
	Sentry  string
	Server  string
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

}
