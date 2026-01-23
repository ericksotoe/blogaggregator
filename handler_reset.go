package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		fmt.Println("There was an error reseting the table")
		return err
	}
	fmt.Println("Successfully reset the table")
	return nil
}
