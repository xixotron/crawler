package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	aTags := doc.Find("body").Find("a")

	urls := []string{}
	for _, tag := range aTags.EachIter() {
		href, ok := tag.Attr("href")
		if !ok {
			continue
		}

		uri, err := baseURL.Parse(href)
		if err != nil {
			return nil, err
		}
		urls = append(urls, uri.String())
	}

	return urls, nil
}
