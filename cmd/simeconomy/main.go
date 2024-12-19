// Package main implements the tool.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "config.yaml", "config file")
	flag.Parse()

	cfg := loadConfig(configFile)

	slog.Info(fmt.Sprintf("loaded config: %s", configFile))

	w := newWorld(cfg)

	slog.Info("world created")

	slog.Info("entering command loop")

	fmt.Printf("\nWelcome to simeconomy %s\n\n", version)
	fmt.Println("Type help and hit ENTER to list available commands.")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		pop, farmers, merchants := w.populationSize()
		fmt.Printf("day: %d\npopulation: %d\nfarmers: %d\nmerchants: %d\nmoney: %d\ngoods: %d\nfood: %d\n",
			w.day,
			pop,
			farmers,
			merchants,
			w.money(),
			w.amountOfGoods(),
			w.amountOfFood())
		fmt.Print("\ncommand> ")
		text, errRead := reader.ReadString('\n')
		if errRead != nil {
			fmt.Printf("input error: %v\n", errRead)
			break
		}
		execute(w, text)
	}

	fmt.Println("exiting")
}
