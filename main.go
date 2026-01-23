package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/ericksotoe/blogaggregator/internal/config"
	"github.com/ericksotoe/blogaggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	// opens a connection to our database
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// creates a database.Queries pointer that we created using sqlc
	dbQueries := database.New(db)

	programState := &state{db: dbQueries, cfg: &cfg}
	cmds := commands{registerredCommands: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)

	cmdArgs := os.Args
	if len(cmdArgs) < 2 {
		fmt.Printf("Please provide 2 arguments in the CLI")
		os.Exit(1)
	}

	cmdName := cmdArgs[1]
	cmdArgs = cmdArgs[2:]
	command := command{name: cmdName, args: cmdArgs}
	err = cmds.run(programState, command)
	if err != nil {
		fmt.Printf("error running the command exited with body %s", err)
		os.Exit(1)
	}
}
