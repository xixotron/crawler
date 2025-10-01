package main

import (
	"net/url"
	"path"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	url.Scheme = ""
	url.Path = url.EscapedPath()
	url.Path = strings.ReplaceAll(url.Path, "//", "/")
	url.Path = path.Clean(url.Path)
	if url.Path[len(url.Path)-1] == '/' {
		url.Path = url.Path[:len(url.Path)-1]
	}
	if url.Path == "." {
		url.Path = ""
	}
	url.Host = strings.ToLower(url.Host)

	return url.Host + url.Path, nil
}
