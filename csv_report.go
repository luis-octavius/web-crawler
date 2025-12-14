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
		return fmt.Errorf("error - writeCSVReport: %v", err)
	}

	writer := csv.NewWriter(file)

	err = writer.Write([]string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"})
	if err != nil {
		return fmt.Errorf("error - writeCSVReport: %v", err)
	}

	for _, page := range pages {
		outLinks := strings.Join(page.OutgoingLinks, ";")
		imgLinks := strings.Join(page.ImageURLs, ";")
		err = writer.Write([]string{page.URL, page.H1, page.FirstParagraph, outLinks, imgLinks})
		if err != nil {
			return fmt.Errorf("error - writeCSVReport: %v", err)
		}
	}

	return nil
}
