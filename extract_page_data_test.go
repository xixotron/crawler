package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_extractPageData(t *testing.T) {
	tests := map[string]struct {
		htmlBody string
		baseURL  string
		expected PageData
	}{
		"h1, fallback p, relative link and image": {
			htmlBody: `<html>
  <body>
    <h1>Test Title</h1>
    <p>This is the first paragraph.</p>
    <a href="/link1">Link 1</a>
    <img src="/image1.jpg" alt="Image 1">
  </body>
</html>`,
			baseURL: "https://blog.boot.dev",
			expected: PageData{
				URL:            "https://blog.boot.dev",
				H1:             "Test Title",
				FirstParagraph: "This is the first paragraph.",
				OutgoingLinks:  []string{"https://blog.boot.dev/link1"},
				ImageURLs:      []string{"https://blog.boot.dev/image1.jpg"},
			},
		},
		"h1, p from main, relative and absolute link and images": {
			htmlBody: `<html>
  <body>
    <h1>Test Title</h1>
    <p>This is the first paragraph.</p>
    <main>
      <h1>Title in main</h1>
      <p>First pargraph of main</p>
      <a href="/link1">Link 1</a>
      <a href="https://google.com">google</a>
      <img src="/image1.jpg" alt="Image 1">
      <img src="https://boot.dev/image2.jpg" alt="Image 2">
    </main>
    <img src="/image3.jpg" alt="Image 3">
    <a href="https://boot.dev/about">about</a>
  </body>
</html>`,
			baseURL: "https://blog.boot.dev",
			expected: PageData{
				URL:            "https://blog.boot.dev",
				H1:             "Test Title",
				FirstParagraph: "First pargraph of main",
				OutgoingLinks: []string{
					"https://blog.boot.dev/link1",
					"https://google.com",
					"https://boot.dev/about",
				},
				ImageURLs: []string{
					"https://blog.boot.dev/image1.jpg",
					"https://boot.dev/image2.jpg",
					"https://blog.boot.dev/image3.jpg",
				},
			},
		},
		"malformed HTML still parsed; absolute link and image": {
			htmlBody: `<html body>
  <h1>Messy</h1>
  <a href="https://other.com/path">Other</a>
  <img src="https://cdn.boot.dev/banner.jpg">
</html body>`,
			baseURL: "https://blog.boot.dev",
			expected: PageData{
				URL:            "https://blog.boot.dev",
				H1:             "Messy",
				FirstParagraph: "", // no paragraph
				OutgoingLinks:  []string{"https://other.com/path"},
				ImageURLs:      []string{"https://cdn.boot.dev/banner.jpg"},
			},
		},
		"invalid base URL → empty link/image slices": {
			htmlBody: `<html>
  <body>
    <h1>Title</h1>
    <p>Paragraph</p>
    <a href="/path">path</a>
    <img src="/logo.png">
  </body>
</html>`,
			baseURL: `:\\invalidBaseURL`,
			expected: PageData{
				URL:            `:\\invalidBaseURL`,
				H1:             "Title",
				FirstParagraph: "Paragraph",
				OutgoingLinks:  nil,
				ImageURLs:      nil,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := extractPageData(tc.htmlBody, tc.baseURL)
			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("Test '%s' FAIL: %v", name, diff)
			}
		})
	}
}
