package main

import (
	"fmt"

	"github.com/ericksotoe/blogaggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cliCommand map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username.\n")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to %s\n", cmd.args[0])
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	// This method runs a given command with the provided state if it exists.
	if s.cfg.DbUrl == "" {
		return fmt.Errorf("The state doesn't have a registered db")
	}

	_, ok := c.cliCommand[cmd.name]
	if !ok {
		return fmt.Errorf("command not found in the commands struct")
	}

	f := c.cliCommand[cmd.name]
	err := f(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	_, ok := c.cliCommand[name]
	if !ok {
		c.cliCommand[name] = f
	}

	// This method registers a new handler function for a command name.
}
