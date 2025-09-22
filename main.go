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
	htmlContent, err := getHTML(rawBaseURL)
	if err != nil {
		fmt.Printf("getHTML: %v\n", err)
	}
	fmt.Println(htmlContent)
}
