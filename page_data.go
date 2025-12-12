package main

import (
	"log"
	"net/url"
)

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(html, pageURL string) PageData {
	URL, err := url.Parse(pageURL)
	if err != nil {
		log.Println("error parsing page URL parameter: ", err)
		return PageData{}
	}

	header := getH1FromHTML(html)

	firstParagraph := getFirstParagraphFromHTML(html)

	URLs, err := getURLsFromHTML(html, URL)
	if err != nil {
		log.Println("error getting URLs from HTML", err)
		return PageData{}
	}

	imagesLinks, err := getImagesFromHTML(html, URL)
	if err != nil {
		log.Println("error getting URLs images from HTML", err)
		return PageData{}
	}

	return PageData{
		URL:            pageURL,
		H1:             header,
		FirstParagraph: firstParagraph,
		OutgoingLinks:  URLs,
		ImageURLs:      imagesLinks,
	}
}
