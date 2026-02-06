package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ericksotoe/blogaggregator/internal/database"
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

	// Get the next feed that should be fetched, ordered by its last fetched time
	fetchedFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("get next feed to fetch: %w", err)
	}

	// Mark this feed as having been fetched now (updates last_fetched_at / updated_at)
	err = s.db.MarkFeedFetched(ctx, fetchedFeed.ID)
	if err != nil {
		return fmt.Errorf("mark feed %s as fetched: %w", fetchedFeed.ID, err)
	}

	// Download and parse the RSS feed from the feed's URL
	rssFeeds, err := fetchFeed(ctx, fetchedFeed.Url)
	if err != nil {
		return fmt.Errorf("Something went wrong when fetching using the url. \n'Error Body': %v\n", err)
	}

	// Loop through each item (post) in the RSS feed
	for _, feedItem := range rssFeeds.Channel.Item {
		// Use the item description, or a default message if it's empty
		desc := feedItem.Description
		if desc == "" {
			desc = "no feed description available"
		}

		// Try to parse the published date using our helper that supports multiple formats
		t, err := parsePubDate(feedItem.PubDate)
		published := sql.NullTime{}

		if err == nil {
			// If parsing succeeded, store it as a valid SQL time
			published = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		// Insert the post into the database, linked to the current feed
		_, err = s.db.CreatePost(ctx, database.CreatePostParams{
			Title: feedItem.Title,
			Url:   feedItem.Link,
			Description: sql.NullString{
				String: desc,
				Valid:  true,
			},
			PublishedAt: published,
			FeedID:      fetchedFeed.ID,
		})

		// If something goes wrong inserting a post, bail out for now
		if err != nil {
			return fmt.Errorf("create post for feed %s (title: %q, url: %q): %w",
				fetchedFeed.ID,
				feedItem.Title,
				feedItem.Link,
				err,
			)
		}
	}
	// All posts for this feed were processed successfully
	return nil
}

var layouts = []string{
	time.RFC1123Z,
	time.RFC1123,
	time.RFC3339,
}

func parsePubDate(s string) (time.Time, error) {
	var lastErr error
	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

func handlerBrowsePosts(s *state, cmd command, user database.User) error {
	var limit int32 = 2
	if len(cmd.args) >= 1 {
		num, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("Error changing string int to integer %w", err)
		}
		limit = int32(num)
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit}

	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error getting posts for logged in user: %w\n", err)
	}
	for _, post := range posts {
		printPost(post)
	}

	return nil
}

func printPost(p database.GetPostsForUserRow) {
	fmt.Printf("Feed: %s\n", p.FeedName)
	fmt.Printf("Title: %s\n", p.Title)
	fmt.Printf("URL: %s\n", p.Url)

	if p.Description.Valid {
		fmt.Printf("Description: %s\n", p.Description.String)
	}

	if p.PublishedAt.Valid {
		fmt.Printf("Published At: %s\n", p.PublishedAt.Time.Format(time.RFC3339))
	}

	fmt.Println("-----")
}
