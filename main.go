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

	config, err := defaultConfig(rawBaseURL, 1)
	if err != nil {
		fmt.Printf("Cound't build config: %v", err)
		os.Exit(1)
	}
	config.wg.Go(func() {
		config.crawlPage(rawBaseURL)
	})

	config.wg.Wait()

	fmt.Printf("Crawl of %s completed\n", rawBaseURL)

	fmt.Println("Results:")
	for normalizedURL, pageData := range config.pages {
		fmt.Printf(" - %s linked: times %d\n", normalizedURL, pageData.Visits)
	}
}
