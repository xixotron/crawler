package main

import (
	"log"
	"net/url"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	if cfg.reachedMaxPages() {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Printf("error crawlPage: Couldn't parse currentURL %q: %v", rawCurrentURL, err)
		return
	}

	if cfg.baseURL.Hostname() != currentURL.Hostname() {
		log.Printf("Hostname %v doesn't match current url %v", currentURL, cfg.baseURL)
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Printf("error crawlPage: Couldn't normalize currentURL %q: %v", rawCurrentURL, err)
		return
	}
	if !cfg.addPageVisit(normalizedURL) {
		return
	}

	log.Printf("Crawling: %q, normalized as: %q", currentURL.String(), normalizedURL)

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		log.Printf("error getHTML: %v", err)
		return
	}

	pageData := extractPageData(html, rawCurrentURL)
	cfg.setPageData(normalizedURL, pageData)

	for _, nextURL := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(nextURL)
	}
}
