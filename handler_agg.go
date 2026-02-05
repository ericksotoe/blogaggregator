package main

import (
	"fmt"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Invalid arguments provided to the agg command, please add a time ex: 1s, 1m, 1h")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Something went wrong when changing the string time into ticker time")
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}
