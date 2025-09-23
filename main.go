package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(os.Args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	rawBaseURL := os.Args[1]
	fmt.Printf("starting crawl of: %s...\n", rawBaseURL)

	pages := make(map[string]int)
	crawlPage(rawBaseURL, rawBaseURL, pages)
	fmt.Printf("Crawl of %s completed\n", rawBaseURL)

	fmt.Println("Results:")
	for page, count := range pages {
		fmt.Printf(" - %s linked: %d times\n", page, count)
	}
}
