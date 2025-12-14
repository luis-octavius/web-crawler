package main

import (
	"fmt"
	"log"
	"net/url"
)

// crawlPage recursively crawls a current url
// it sends a signal to the concurrencyControl channel
// for controlling the goroutine executed from main
//
// returns an error if:
// - the domain of the current URL is different than the base URL
// - normalizeURL returns an error
// - isFirst returns false
func (cfg *config) crawlPage(rawCurrentURL string) {
	// sends a signal into the concurrencyControl channel
	cfg.concurrencyControl <- struct{}{}

	// when the function returns receives the signal sent before
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	if cfg.pagesLen() >= cfg.maxPages {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("error - crawlPage: couldn't parse URL '%s': %v", rawCurrentURL, err)
	}

	if currentURL.Hostname() != cfg.baseURL.Hostname() {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}

	isFirst := cfg.addPageVisit(normalizedURL)
	if !isFirst {
		return
	}

	// fmt.Printf("crawling %s\n", rawCurrentURL)

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error - getHTML: %v\n", err)
		return
	}

	pageData := extractPageData(html, rawCurrentURL)
	cfg.setPageData(normalizedURL, pageData)

	for _, URL := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(URL)
	}
}
