package main

import "testing"

func Test_getH1FromHTML(t *testing.T) {
	tests := map[string]struct {
		inputURL string
		expected string
	}{
		"h1 content from basic HTML": {
			inputURL: "<html><body><h1>Test Title</h1></body></html>",
			expected: "Test Title",
		},

		"only first h1 content when multiple present": {
			inputURL: `<html>
  <body>
    <h1>Welcome to Boot.dev</h1>
    <h1>We have some titles for you</h1>
    <main>
      <h1>First you can:</h1>
      <p>Learn to code by building real projects.</p>
      <p>This is the second paragraph.</p>
    </main>
  </body>
</html>`,
			expected: "Welcome to Boot.dev",
		},

		"empty string when no h1 element present": {
			inputURL: "<html><body><p>not h1 an h1 element</p><h2>also not h1</h2></body></html>",
			expected: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := getH1FromHTML(tc.inputURL)
			if actual != tc.expected {
				t.Errorf("Test '%s' FAIL: expected string: %q, actual: %q", name, tc.expected, actual)
			}
		})
	}
}

func Test_getFirstParagraphFromHTML(t *testing.T) {

	tests := map[string]struct {
		inputURL string
		expected string
	}{
		"p content from basic HTML": {
			inputURL: "<html><body><p>Test Paragraph</p></body></html>",
			expected: "Test Paragraph",
		},

		"only first p content when multiple present": {
			inputURL: `<html>
  <body>
    <h1>Welcome to Boot.dev</h1>
    <p>We have some paragraphs for you</p>
    <p>Learn to code by building real projects.</p>
    <p>This is the third paragraph.</p>
  </body>
</html>`,
			expected: "We have some paragraphs for you",
		},

		"first p from main when present": {
			inputURL: `<html>
  <body>
    <h1>Welcome to Boot.dev</h1>
    <p>We have some paragraphs for you</p>
      <main>
        <p>Learn to code by building real projects.</p>
        <p>This is the third paragraph.</p>
      </main>
  </body>
</html>`,
			expected: "Learn to code by building real projects.",
		},

		"empty string when no p element present": {
			inputURL: "<html><body><h1>not a p element</h1><h2>also not p</h2></body></html>",
			expected: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tc.inputURL)
			if actual != tc.expected {
				t.Errorf("Test '%s' FAIL: expected string: %q, actual: %q", name, tc.expected, actual)
			}
		})
	}
}
