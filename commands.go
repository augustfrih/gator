package main

import (
	"fmt"
)

type command struct {
	name      string
	arguments []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	commandToRun, ok := c.registeredCommands[cmd.name]
	if !ok {
		return fmt.Errorf("Tried to run nonexistent command: %v", cmd.name)
	}
	return commandToRun(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) error {
	c.registeredCommands[name] = f
	return nil
}
