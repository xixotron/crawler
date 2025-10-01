package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func writeCSVReport(pages map[string]PageData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	writer := csv.NewWriter(file)

	err = writer.Write([]string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"})
	if err != nil {
		return fmt.Errorf("erro writing header line: %w", err)
	}
	for _, pageData := range pages {
		err := writer.Write([]string{
			pageData.URL,
			pageData.H1,
			pageData.FirstParagraph,
			strings.Join(pageData.OutgoingLinks, ";"),
			strings.Join(pageData.ImageURLs, ";"),
		})
		if err != nil {
			return fmt.Errorf("error writing entry: %w", err)
		}
	}

	return nil
}
