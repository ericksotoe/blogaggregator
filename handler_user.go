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

	_, err := s.db.CreateUser(context.Background(), userToCreate)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}

	// fmt.Printf("User: %s was created in the db with the following information\n", name)
	// fmt.Printf("ID: %d\nCreated At: %v\nUpdated At: %v\nName: %s", user.ID, user.CreatedAt, user.UpdatedAt, user.Name)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("the login command expects a single argument, the username.\n")
	}
	name := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		fmt.Printf("user exists, body of error %s\n", err)
		os.Exit(1)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to %s\n", name)
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("No users found in the db: %w\n", err)
	}

	if len(users) == 0 {
		fmt.Println("No users found in the database. Try creating one with 'register <name>'")
		return nil
	}

	currentlyLoggedUser := s.cfg.Username

	for _, user := range users {
		if user == currentlyLoggedUser && currentlyLoggedUser != "" {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}
