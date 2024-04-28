package main

import (
	"flag"
	"log"

	"github.com/aedobrynin/soa-hw/statistics/internal/app"
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

	a, err := app.New(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := a.Serve(); err != nil {
		log.Fatal(err.Error())
	}
}
