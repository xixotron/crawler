package main

import (
	"fmt"
	"log"
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]*PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func defaultConfig(rawBaseURL string, concurrntProcesses int) (*config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("error Couldn't parse baseURL %q: %v", rawBaseURL, err)
	}
	return &config{
		pages:              make(map[string]*PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, concurrntProcesses),
		wg:                 &sync.WaitGroup{},
	}, nil

}

func (cfg *config) crawlPage(rawCurrentURL string) {
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

	log.Printf("crawling: %s", rawCurrentURL)

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		log.Printf("error getHTML: %v", err)
		return
	}

	pageData := extractPageData(html, rawCurrentURL)
	pageData.Visits = 1
	cfg.mu.Lock()
	cfg.pages[normalizedURL] = &pageData
	cfg.mu.Unlock()

	for _, nextURL := range pageData.OutgoingLinks {
		cfg.concurrencyControl <- struct{}{}
		cfg.wg.Go(func() {
			cfg.crawlPage(nextURL)
		})
		<-cfg.concurrencyControl
	}
}

func (cfg *config) addPageVisit(normalizedURL string) bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, ok := cfg.pages[normalizedURL]; ok {
		cfg.pages[normalizedURL].Visits += 1
		return false
	}

	cfg.pages[normalizedURL] = &PageData{}
	return true
}
