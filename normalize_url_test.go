package main

import "testing"

func Test_normalizeURL(t *testing.T) {
	tests := map[string]struct {
		inputURL string
		expected string
	}{
		"remove https scheme": {
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		"remove http scheme": {
			inputURL: "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		"remove final separator": {
			inputURL: "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		"remove duplicate slashes": {
			inputURL: "http://blog.boot.dev//path//",
			expected: "blog.boot.dev/path",
		},
		"remove dot segments": {
			inputURL: "http://blog.boot.dev/./path/../",
			expected: "blog.boot.dev",
		},
		"lowercase host part": {
			inputURL: "http://BLOG.BOOT.DEV/path",
			expected: "blog.boot.dev/path",
		},
		"base url without path": {
			inputURL: "https://wagslane.dev",
			expected: "wagslane.dev",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test '%s' FAIL: unexpected error: %v", name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test '%s' FAIL: expected URL: %q, actual: %q", name, tc.expected, actual)
			}
		})
	}
}
