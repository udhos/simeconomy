package main

import (
	"fmt"
	"os"
	"strings"
)

type command struct {
	name        string
	exec        func(w *world)
	description string
}

var commandTable []command

func init() {
	commandTable = []command{
		{"?", cmdHelp, "List commands"},
		{"help", cmdHelp, "List commands"},
		{"quit", cmdQuit, "Exit simulator"},
		{"run", cmdRun, "Run simulation step"},
	}
}

func execute(w *world, cmd string) {
	cmd = strings.TrimSpace(cmd)
	for _, c := range commandTable {
		if strings.HasPrefix(c.name, cmd) {
			c.exec(w)
			return
		}
	}
	fmt.Printf("command not found: '%s'", cmd)
}

func cmdHelp(_ *world) {
	for _, c := range commandTable {
		fmt.Printf("%s - %s\n", c.name, c.description)
	}
}

func cmdQuit(_ *world) {
	fmt.Println("quitting")
	os.Exit(0)
}

func cmdRun(w *world) {
	w.day++
	fmt.Println("simulation step executed")
}
