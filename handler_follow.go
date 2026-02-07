package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ericksotoe/blogaggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("the follow command expects a single argument, the url.\n")
	}

	feedURL := cmd.args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Error getting the feed from the given url.")
	}

	feedFollowToCreate := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowToCreate)
	if err != nil {
		// check for "duplicate key"
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			fmt.Println("You are already following this feed.")
			return nil
		}
		return fmt.Errorf("error adding to feed_follows (user %s, feed %s): %w", user.ID, feed.ID, err)
	}
	for _, feed := range feedFollow {
		fmt.Printf("feed name: %s\ncurrent user: %s", feed.FeedName, feed.UserName)
	}

	return nil
}

func handlerFeedFollowsForUser(s *state, cmd command, user database.User) error {

	feedsFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Error getting the feed following using the user ID")
	}

	for _, feed := range feedsFollows {
		fmt.Printf("User: %s is following %s\n", feed.UserName, feed.FeedName)
	}

	return nil
}

func handlerFeedDelete(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Invalid argument provided to unfollow command, the url must be provided")
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])

	feedToDelete := database.DeleteUsingBothIDSParams{
		FeedID: feed.ID,
		UserID: user.ID,
	}

	_, err = s.db.DeleteUsingBothIDS(context.Background(), feedToDelete)
	if err != nil {
		return fmt.Errorf("Error deleteing the feed using both feed and user id")
	}

	fmt.Printf("%s Unfollowed successfully!\n", feed.Name)
	return nil
}
