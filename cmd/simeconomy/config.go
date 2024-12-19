package main

import (
	"fmt"
	"log/slog"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type config struct {
	InitialPopulation int `yaml:"initial_population"`
	InitialMoney      int `yaml:"initial_money"`
}

func loadConfig(configFile string) config {
	input, errRead := os.ReadFile(configFile)
	if errRead != nil {
		slog.Error(fmt.Sprintf("FATAL: loadConfig: file error: %v", errRead))
		os.Exit(1)
	}
	var cfg config
	errConf := yaml.Unmarshal(input, &cfg)
	if errConf != nil {
		slog.Error(fmt.Sprintf("FATAL: loadConfig: yaml error: %v", errRead))
		os.Exit(1)
	}
	return cfg
}
