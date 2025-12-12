package main

import (
	"fmt"
	"log"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {

	isSameDomain, err := checkDomain(rawBaseURL, rawCurrentURL)
	if !isSameDomain {
		log.Printf("error: %v\n", err)
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}

	page, ok := pages[normalizedURL]
	if ok {
		page += 1
		return
	} else {
		pages[normalizedURL] = 1
	}

	html, err := getHTML(rawCurrentURL)
	fmt.Println("HTML: ", html)

	pageData := extractPageData(html, rawCurrentURL)
	fmt.Println("Page Data: ", pageData)

	for _, URL := range pageData.OutgoingLinks {
		crawlPage(rawBaseURL, URL, pages)
	}
}

func checkDomain(rawBaseURL, rawCurrentURL string) (bool, error) {
	parsedBaseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return false, fmt.Errorf("error parsing base URL: %v", err)
	}

	parsedCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return false, fmt.Errorf("error parsing current URL: %v", err)
	}

	baseDomain := parsedBaseURL.Hostname()
	currentDomain := parsedCurrentURL.Hostname()

	if baseDomain != currentDomain {
		return false, fmt.Errorf("base URL with domain %v has a different domain than current URL %v", baseDomain, currentDomain)
	}

	return true, nil
}
