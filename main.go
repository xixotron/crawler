package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	const defaultMaxConcurency = 5
	const defaultMaxPages = 20

	maxConcurency := defaultMaxConcurency
	maxPages := defaultMaxPages

	if len(os.Args) < 2 {
		fmt.Println("no URL provided")
		print_usage()
		os.Exit(1)
	}
	if len(os.Args) > 4 {
		fmt.Println("too many arguments provided")
		print_usage()
		os.Exit(1)
	}

	rawBaseURL := os.Args[1]

	if len(os.Args) > 2 {
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("invalid maxConcurency value: %v\n", err)
			print_usage()
			os.Exit(1)
		}
		maxConcurency = n
	}
	if len(os.Args) > 3 {
		n, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Printf("invalid maxPages value: %v\n", err)
			print_usage()
			os.Exit(1)
		}
		maxPages = n
	}

	config, err := configure(rawBaseURL, maxConcurency, maxPages)
	fmt.Printf("starting crawl of: %q with maxConcurency: %d, maxPages: %d ...\n", config.baseURL.String(), cap(config.concurrencyControl), config.maxPages)
	if err != nil {
		fmt.Printf("Cound't build config: %v", err)
		os.Exit(1)
	}
	config.wg.Add(1)
	go config.crawlPage(rawBaseURL)

	config.wg.Wait()

	fmt.Printf("Crawl of %s completed\n", rawBaseURL)

	fmt.Println("Results:")
	for normalizedURL := range config.pages {
		fmt.Printf(" - %s visited\n", normalizedURL)
	}
}

func print_usage() {
	binaryParts := strings.Split(os.Args[0], "/")
	binaryName := binaryParts[len(binaryParts)-1]
	fmt.Println("usage:")
	fmt.Printf("%s URL [maxConcurrency [maxPages]]\n", binaryName)
	fmt.Println(" URL:               URI for site to crawl")
	fmt.Println(" maxConcurrency:    maximum number of parallel requests to run at a time")
	fmt.Println(" maxPages:          maximum number of pages to crawl")
}
