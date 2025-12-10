package main

import (
	"fmt"
	"net/url"
)

func normalizeURL(s string) (string, error) {
	url, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("error parsing url: %v", err)
	}

	cleanedURL := url.Host + "/" + url.Path
	return cleanedURL, nil
}
