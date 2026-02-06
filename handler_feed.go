package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ericksotoe/blogaggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("the addfeed handler expects two arguments, the name and url.\n")
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]
	feedToAdd := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedToAdd)
	if err != nil {
		fmt.Printf("error creating the feed with following passed in params\nName: %s\nURL: %s\nError body: %s", feedName, feedURL, err)
		os.Exit(1)
	}

	feedFollowToCreate := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowToCreate)
	if err != nil {
		return fmt.Errorf("Error adding to the feed_follows table")
	}
	return nil
}

func handlerGetFeed(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Ran into an error when getting the feed from the feeds db/ body:%s", err)
	}

	if len(feeds) == 0 {
		fmt.Println("The database doesn't have any feeds")
		return nil
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
