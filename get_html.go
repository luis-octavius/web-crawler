package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var Client = &http.Client{
	Timeout: 30 * time.Second,
}

// getHTML busca o HTML de uma URL com retry automático e rate limiting
func getHTML(rawURL string) (string, error) {
	retryConfig := DefaultRetryConfig()

	var html string
	err := RetryWithBackoff(retryConfig, func() error {
		var fetchErr error
		html, fetchErr = fetchHTML(rawURL)
		return fetchErr
	})

	return html, err
}

// fetchHTML faz a requisição HTTP real
func fetchHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição: %v", err)
	}

	// Headers mais realistas para evitar bloqueios
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	res, err := Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao fazer requisição: %v", err)
	}
	defer res.Body.Close()

	// Tratamento específico para diferentes códigos de status
	switch {
	case res.StatusCode == 429:
		return "", fmt.Errorf("rate limit atingido (429) - aguardando retry")
	case res.StatusCode >= 500:
		return "", fmt.Errorf("erro no servidor (%d) - aguardando retry", res.StatusCode)
	case res.StatusCode >= 400:
		return "", fmt.Errorf("erro do cliente (%d) - não retentável", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("content type não é text/html: %v", contentType)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler corpo da resposta: %v", err)
	}

	return string(body), nil
}

// getHTMLWithContext busca HTML com contexto para cancelamento
func getHTMLWithContext(ctx context.Context, rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição: %v", err)
	}

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	res, err := Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao fazer requisição: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("status code: %v", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %v", err)
	}

	return string(body), nil
}
