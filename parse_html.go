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

	log.Println("Header: ", header)
	return header
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

	log.Println("First p: ", firstParagraph)
	return firstParagraph
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	reader := strings.NewReader(htmlBody)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return []string{""}, fmt.Errorf("error creating the document from reader: %v", err)
	}

	parsedURLs := []string{}
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		parsedURLs = append(parsedURLs, url)
	})

	for _, url := range parsedURLs {
		fmt.Println("URL: ", url)
	}

	return parsedURLs, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	reader := strings.NewReader(htmlBody)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return []string{""}, fmt.Errorf("error creating the document from reader: %v", err)
	}

	parsedURLs := []string{}
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("src")
		absoluteURL := baseURL.String() + url
		parsedURLs = append(parsedURLs, absoluteURL)
	})

	return parsedURLs, nil
}
