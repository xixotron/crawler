package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	aTags := doc.Find("body").Find("img")

	urls := []string{}
	for _, tag := range aTags.EachIter() {
		src, ok := tag.Attr("src")
		if !ok {
			continue
		}

		uri, err := baseURL.Parse(src)
		if err != nil {
			return nil, err
		}
		urls = append(urls, uri.String())
	}

	return urls, nil
}
