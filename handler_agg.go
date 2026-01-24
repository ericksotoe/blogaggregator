package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Printf("Title: %s\nLink: %s\nDescription: %s\n", feed.Channel.Title, feed.Channel.Link, feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		fmt.Printf("Title: %s\nLink: %s\nDescription: %s\n", item.Title, item.Link, item.Description)
	}
	return nil
}
