package main

import (
	"fmt"

	"github.com/ericksotoe/blogaggregator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s is the url again\n", c.DbUrl)
	}
}
