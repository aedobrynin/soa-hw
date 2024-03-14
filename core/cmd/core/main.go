package main

import (
	"flag"
	"log"

	"core/internal/app"
)

func getConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "c", "./.config.yaml", "path to config file")
	flag.Parse()

	return configPath
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
