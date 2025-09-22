package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func getHTML(rawURL string) (string, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", "BootCrawler/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("error response status: %v, %q", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	if contentType := resp.Header.Get("Content-Type"); !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("error response Content-Type: %s ", contentType)
	}

	htmlBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}
	return string(htmlBodyBytes), nil
}
