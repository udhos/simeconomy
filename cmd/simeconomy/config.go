package main

import (
	"fmt"
	"log/slog"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type config struct {
	InitialPopulation         int     `yaml:"initial_population"`
	InitialMoney              int     `yaml:"initial_money"`
	InitialFarmers            int     `yaml:"initial_farmers"`
	InitialMerchants          int     `yaml:"initial_merchants"`
	FruitDuration             int     `yaml:"fruit_duration"`
	DailyMeals                int     `yaml:"daily_meals"`
	FarmerProduction          int     `yaml:"farmer_production"`
	FarmerInitialPrice        int     `yaml:"farmer_initial_price"`
	MerchantInitialSellFactor float32 `yaml:"merchant_initial_sell_factor"`
	MerchantCapacity          int     `yaml:"merchant_capacity"`
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
