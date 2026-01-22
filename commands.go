package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	registerredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if s.cfg.DbUrl == "" {
		return fmt.Errorf("The state doesn't have a registered db")
	}

	f, ok := c.registerredCommands[cmd.name]
	if !ok {
		return fmt.Errorf("command not found in the commands struct")
	}

	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	_, ok := c.registerredCommands[name]
	if !ok {
		c.registerredCommands[name] = f
	}
}
