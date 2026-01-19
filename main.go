package main

import (
	"fmt"

	"github.com/ericksotoe/blogaggregator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	err = c.SetUser("erick")
	if err != nil {
		fmt.Println(err)
	}

	c, err = config.Read()
	fmt.Printf("%s\n%s\n", c.DbUrl, c.Username)
}
