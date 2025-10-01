package main

import "net/url"

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(html, pageURL string) PageData {
	h1 := getH1FromHTML(html)
	firstParagraph := getFirstParagraphFromHTML(html)

	URL, err := url.Parse(pageURL)
	if err != nil {
		return PageData{
			URL:            pageURL,
			H1:             h1,
			FirstParagraph: firstParagraph,
			OutgoingLinks:  nil,
			ImageURLs:      nil,
		}
	}

	outgoingLinks, err := getURLsFromHTML(html, URL)
	if err != nil {
		outgoingLinks = nil
	}

	imageURLs, err := getImagesFromHTML(html, URL)
	if err != nil {
		imageURLs = nil
	}

	return PageData{
		URL:            pageURL,
		H1:             h1,
		FirstParagraph: firstParagraph,
		OutgoingLinks:  outgoingLinks,
		ImageURLs:      imageURLs,
	}
}
