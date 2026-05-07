package main

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"
)

type config struct {
	pages              map[string]PageData
	maxPages           int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	rateLimiter        *RateLimiter
	proxyRotator       *ProxyRotator
	requestDelay       time.Duration
	ctx                context.Context
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	_, visited := cfg.pages[normalizedURL]
	if visited {
		return false
	}

	cfg.pages[normalizedURL] = PageData{URL: normalizedURL}
	return true
}

func (cfg *config) setPageData(normalizedURL string, data PageData) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.pages[normalizedURL] = data
}

func (cfg *config) pagesLen() int {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages)
}

func configure(rawBaseURL string, maxPages, maxConcurrency int) (*config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	// Configuração de rate limiting: 2 requisições por segundo, burst de 5
	rateLimiter := NewRateLimiter(2.0, 5)

	return &config{
		pages:              make(map[string]PageData),
		maxPages:           maxPages,
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		rateLimiter:        rateLimiter,
		proxyRotator:       nil, // Pode ser configurado depois
		requestDelay:       500 * time.Millisecond,
		ctx:                context.Background(),
	}, nil
}

// ConfigureWithProxies configura o crawler com rotação de proxies
func ConfigureWithProxies(rawBaseURL string, maxPages, maxConcurrency int, proxyURLs []string) (*config, error) {
	cfg, err := configure(rawBaseURL, maxPages, maxConcurrency)
	if err != nil {
		return nil, err
	}

	if len(proxyURLs) > 0 {
		proxyRotator, err := NewProxyRotator(proxyURLs, true)
		if err != nil {
			return nil, fmt.Errorf("erro ao configurar proxies: %v", err)
		}
		cfg.proxyRotator = proxyRotator
	}

	return cfg, nil
}
