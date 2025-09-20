package main

import (
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_getURLsFromHTML(t *testing.T) {
	tests := map[string]struct {
		htmlBody    string
		baseURL     string
		expected    []string
		expectedErr bool
	}{
		"url from absolute link": {
			htmlBody:    `<html><body><a href="https://blog.boot.dev"><span>Boot.dev</span></a></body></html>`,
			baseURL:     "https://boot.dev",
			expected:    []string{"https://blog.boot.dev"},
			expectedErr: false,
		},
		"url from relative link": {
			htmlBody:    `<html><body><a href="/posts/"><span>Boot.dev</span></a></body></html>`,
			baseURL:     "https://blog.boot.dev",
			expected:    []string{"https://blog.boot.dev/posts/"},
			expectedErr: false,
		},
		"error from malformed link": {
			htmlBody:    `<html><body><a href=":\\invalidURL"><span>Boot.dev</span></a></body></html>`,
			baseURL:     "https://blog.boot.dev",
			expected:    nil,
			expectedErr: true,
		},
		"multiple urls from multiple a tags": {
			htmlBody: `<html>
  <body>
    <a href="/courses/"><span>Courses</span></a>
    <main>
      <a href="https://blog.boot.dev/posts/"><span>Blog posts!</span></a>
    </main>
  </body>
</html>`,
			baseURL: "https://boot.dev",
			expected: []string{
				"https://boot.dev/courses/",
				"https://blog.boot.dev/posts/",
			},
			expectedErr: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			url, err := url.Parse(tc.baseURL)
			if err != nil {
				t.Errorf("Test '%s' Error parsing baseURL: %q error: %v", name, tc.baseURL, err)
				return
			}

			actual, err := getURLsFromHTML(tc.htmlBody, url)
			if (err != nil) != tc.expectedErr {
				t.Errorf("Test '%s' FAIL: expected error: %v, got: %v", name, tc.expectedErr, err)
			}
			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("Test '%s' FAIL: %v", name, diff)
			}
		})
	}
}
