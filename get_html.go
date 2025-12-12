package main

import (
	"fmt"
	"io"
	"net/http"
)

var Client = &http.Client{}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("error doing the request with the provived URL: %v", err)
	}

	req.Header.Set("User-Agent", "BootCrawler/1.0")

	res, err := Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error getting the response from the created request: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("status code is not successful: %v", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		return "", fmt.Errorf("content type is not text/html: %v", contentType)
	}

	body, err := io.ReadAll(res.Body)

	return string(body), nil
}
