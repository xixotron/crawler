package main

import (
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_getImagesFromHTML(t *testing.T) {
	tests := map[string]struct {
		htmlBody    string
		baseURL     string
		expected    []string
		expectedErr bool
	}{
		"image from absolute img src": {
			htmlBody:    `<html><body><img src="https://blog.boot.dev/boots.png" /></body></html>`,
			baseURL:     "https://boot.dev",
			expected:    []string{"https://blog.boot.dev/boots.png"},
			expectedErr: false,
		},
		"image from relative img src": {
			htmlBody:    `<html><body><img src="/boots.png" /></body></html>`,
			baseURL:     "https://blog.boot.dev",
			expected:    []string{"https://blog.boot.dev/boots.png"},
			expectedErr: false,
		},
		"error from malformed img src": {
			htmlBody:    `<html><body><img src=":\\invalidURL" /></body></html>`,
			baseURL:     "https://blog.boot.dev",
			expected:    nil,
			expectedErr: true,
		},
		"multiple images from multiple img tags": {
			htmlBody: `<html>
  <body>
    <img src="https://blog.boot.dev/boots.png" />
    <main>
      <img src="/boots.png" />
    </main>
  </body>
</html>`,
			baseURL: "https://boot.dev",
			expected: []string{
				"https://blog.boot.dev/boots.png",
				"https://boot.dev/boots.png",
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

			actual, err := getImagesFromHTML(tc.htmlBody, url)
			if (err != nil) != tc.expectedErr {
				t.Errorf("Test '%s' FAIL: expected error: %v, got: %v", name, tc.expectedErr, err)
			}
			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("Test '%s' FAIL: %v", name, diff)
			}
		})
	}
}
