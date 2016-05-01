package dockrouted

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/formatters/logstash"
	"github.com/alexflint/go-arg"
	"github.com/caarlos0/env"
	"net/http"
)

type Args struct {
	Listen   string `arg:"-l,help: [Host:]Port the address to listen on (:8080)" env:"DOCKROUTE_LISTEN"`
	LogInfo  bool   `arg:"--log-info,help: Log on INFO level (false)" env:"DOCKROUTE_LOG_INFO"`
	LogDebug bool   `arg:"--log-debug,help: Log on DEBUG level (false)" env:"DOCKROUTE_LOG_DEBUG"`
	LogJSON  bool   `arg:"--log-json,help: Log in JSON format (false)" env:"DOCKROUTE_LOG_JSON"`
}

func Main() {
	defer func() {
		if p := recover(); p != nil {
			log.WithFields(log.Fields{"panic": p}).Info("Got panic while startup")
			os.Exit(1)
		}
	}()

	args := loadArgs()
	if args.LogJSON {
		log.SetFormatter(&logstash.LogstashFormatter{TimestampFormat: time.RFC3339Nano})
	} else {
		//log.SetFormatter(&log.TextFormatter{DisableColors: true})
	}
	if args.LogInfo {
		log.SetLevel(log.InfoLevel)
	}
	if args.LogDebug {
		log.SetLevel(log.DebugLevel)
	}

	StartupService(args)
	waitForTermination(func() {
	})
}

func StartupService(args Args) {
	backend := NewDockerDiscovery()
	log.WithFields(log.Fields{"addr": args.Listen}).Info("start dockrouted")
	proxy := NewProxy(backend)
	err := http.ListenAndServe(args.Listen, proxy)
	if err != nil {
		log.WithError(err).Error("failed to listen")
	}
}

func loadArgs() Args {
	args := Args{
		Listen: ":80",
	}

	env.Parse(&args)
	arg.MustParse(&args)
	return args
}

func waitForTermination(callback func()) {
	sigc := make(chan os.Signal)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	log.Infof("Got singal %q: exit greacefully now", <-sigc)
	callback()
	log.Info("exit now")
	os.Exit(0)
}
