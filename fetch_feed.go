package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	request, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "gator")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	feed := RSSFeed{}
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	err = decodeEscapedHTML(&feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}

func decodeEscapedHTML(feed *RSSFeed) error {
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)

	for _, item := range feed.Channel.Item {
		item.Description = html.UnescapeString(item.Description)
		item.Title = html.UnescapeString(item.Title)
	}

	return nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	fetchedFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("Something went wrong when fetching the next feed from oldest to newest %v", err)
	}

	err = s.db.MarkFeedFetched(ctx, fetchedFeed.ID)
	if err != nil {
		return fmt.Errorf("Something went wrong when updating the fetched and updated at column %v", err)
	}

	rssFeeds, err := fetchFeed(ctx, fetchedFeed.Url)
	if err != nil {
		return fmt.Errorf("Something went wrong when fetching using the url %v", err)
	}
	for _, feedItem := range rssFeeds.Channel.Item {
		fmt.Println(feedItem.Title)
	}
	return nil

}
