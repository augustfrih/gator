package main

import (
	"fmt"
	"os"

	"github.com/augustfrih/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	var sta state
	sta.cfg = &cfg

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	err = cmds.register("login", handlerLogin)
	if err != nil {
		fmt.Println("Couldnt register handlerLogin with error:", err)
		os.Exit(1)
	}

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Not enough arguments were provided")
		os.Exit(1)
	}
	var cmd command

	cmd.name, cmd.arguments = args[1], args[2:]

	err = cmds.run(&sta, cmd)
	if err != nil {
		fmt.Printf("Couldnt run command %s. Error: %s", cmd.name, err)
		os.Exit(1)
	}
}
