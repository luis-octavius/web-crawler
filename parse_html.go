package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getH1FromHTML(html string) string {
	reader := strings.NewReader(html)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return ""
	}

	header := doc.Find("h1").Text()
	if header == "" {
		return header
	}

	return strings.TrimSpace(header)
}

func getFirstParagraphFromHTML(html string) string {
	reader := strings.NewReader(html)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return ""
	}

	// check if the document has a main tag
	// if main tag exists, find the first paragraph inside main
	// otherwise just finds the first paragraph from html body
	isMain := doc.Find("main").Text()
	var firstParagraph string
	if isMain != "" {
		firstParagraph = doc.Find("main").Find("p").First().Text()
		return firstParagraph
	}

	firstParagraph = doc.Find("p").First().Text()

	return strings.TrimSpace(firstParagraph)
}

// getURLsFromHTML receives a HTML body and a base URL to reconstruct
// absolute URLs, after that, it returns an array of all the reconstructed URLs
//
// returns an error if the creation of the document with goquery fails
func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	reader := strings.NewReader(htmlBody)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return []string{""}, fmt.Errorf("error creating the document from reader: %v", err)
	}

	parsedURLs := []string{}
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")

		// checks if the actual URL is equal to baseURL
		if url != baseURL.String() {
			absoluteURL := baseURL.JoinPath(url)
			parsedURLs = append(parsedURLs, absoluteURL.String())
		} else {
			parsedURLs = append(parsedURLs, baseURL.String())
		}
	})

	return parsedURLs, nil
}

// getImagesFromHTML receives a HTML body and a base URL to reconstruct
// absolute URLs from the images, after that, it returns an array of all
// the reconstructed links
//
// returns an error if the creation of the document with goquery fails
func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	reader := strings.NewReader(htmlBody)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return []string{""}, fmt.Errorf("error creating the document from reader: %v", err)
	}

	parsedURLs := []string{}
	// find the img URL then structure the absolute URL of the img
	// for each URL found, then append it to the struct
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("src")
		absoluteURL := baseURL.JoinPath(url)
		log.Println("absolute URL: ", absoluteURL)
		parsedURLs = append(parsedURLs, absoluteURL.String())
	})

	return parsedURLs, nil
}
