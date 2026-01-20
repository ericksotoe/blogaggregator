package main

import (
	"fmt"
	"os"

	"github.com/ericksotoe/blogaggregator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	s := &state{cfg: &c}
	commands := commands{cliCommand: make(map[string]func(*state, command) error)}
	commands.register("login", handlerLogin)
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Please provide 2 arguments in the CLI")
		os.Exit(1)
	}

	commandName := args[1]
	args = args[2:]
	command := command{name: commandName, args: args}
	err = commands.run(s, command)
	if err != nil {
		fmt.Printf("error running the command exited with body %s", err)
		os.Exit(1)
	}
}
