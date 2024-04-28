package main

import (
	"flag"
	"log"

	"github.com/aedobrynin/soa-hw/posts/internal/app"
	"github.com/aedobrynin/soa-hw/posts/internal/logger"
)

func getConfigPath() string {
	var flagConfigPath string

	flag.StringVar(&flagConfigPath, "c", "./.config.yaml", "path to config file")
	flag.Parse()

	return flagConfigPath
}

func main() {
	config, err := app.NewConfig(getConfigPath())
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger.GetLogger(config.App.Debug)
	if err != nil {
		log.Fatal(err)
	}

	a, err := app.New(logger, config)
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	if err := a.Serve(); err != nil {
		logger.Sugar().Fatal(err)
	}
}
