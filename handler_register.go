package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ericksotoe/blogaggregator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("the login handler expects a single argument, the username.\n")
	}

	name := cmd.args[0]

	userToCreate := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	user, err := s.db.CreateUser(context.Background(), userToCreate)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Printf("User: %s was created in the db with the following information\n", name)
	fmt.Printf("ID: %d\nCreated At: %v\nUpdated At: %v\nName: %s", user.ID, user.CreatedAt, user.UpdatedAt, user.Name)
	return nil
}
