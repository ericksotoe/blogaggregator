package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ericksotoe/blogaggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("the addfeed handler expects two arguments, the name and url.\n")
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]
	user, err := s.db.GetUser(context.Background(), s.cfg.Username)
	if err != nil {
		fmt.Printf("user exists, body of error %s\n", err)
	}
	feedToAdd := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}

	_, err = s.db.CreateFeed(context.Background(), feedToAdd)
	if err != nil {
		fmt.Printf("error creating the feed with following passed in params\nName: %s\nURL: %s\nError body: %s", feedName, feedURL, err)
		os.Exit(1)
	}

	// add a printing helper for the different structs and print out this struct

	return nil
}

func handlerGetFeed(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Ran into an error when getting the feed from the feeds db/ body:%s", err)
	}

	for _, feed := range feeds {
		fmt.Printf(" - %s\n", feed.Name)
		fmt.Printf(" - %s\n", feed.Url)
		feedOwner, err := s.db.GetUserFromID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Ran into an error when getting the user from the users db using their id/ body:%s", err)
		}
		fmt.Printf(" - %s\n\n", feedOwner)
	}
	return nil
}
