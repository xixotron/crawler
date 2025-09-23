package main

import (
	"log"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		log.Printf("error crawlPage: Couldn't parse baseURL %q: %v", rawBaseURL, err)
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Printf("error crawlPage: Couldn't parse currentURL %q: %v", rawCurrentURL, err)
		return
	}

	if baseURL.Hostname() != currentURL.Hostname() {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Printf("error crawlPage: Couldn't normalize currentURL %q: %v", rawCurrentURL, err)
		return
	}

	if _, ok := pages[normalizedURL]; ok {
		pages[normalizedURL] += 1
		return
	} else {
		pages[normalizedURL] = 1
	}

	log.Printf("crawling: %s", rawCurrentURL)

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		log.Printf("error getHTML: %v", err)
		return
	}

	urls, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		log.Printf("error getURLsFromHTML: %v", err)
		return
	}

	for _, nextURL := range urls {
		crawlPage(rawBaseURL, nextURL, pages)
	}

}
