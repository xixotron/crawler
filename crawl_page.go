package main

import (
	"fmt"
	"log"
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]PageData
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
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, concurrntProcesses),
		wg:                 &sync.WaitGroup{},
	}, nil

}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

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

func (cfg *config) setPageData(normalizedURL string, pageData PageData) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.pages[normalizedURL] = pageData
}

func (cfg *config) addPageVisit(normalizedURL string) bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, ok := cfg.pages[normalizedURL]; ok {
		return false
	}
	return true
}
