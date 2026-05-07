package main

import (
	"fmt"
	"log"
	"net/url"
	"time"
)

// crawlPage recursivamente faz crawl de uma URL
// Implementa rate limiting e delay entre requisições para evitar bloqueios
//
// retorna um erro se:
// - o domínio da URL atual é diferente do domínio base
// - normalizeURL retorna um erro
// - isFirst retorna false
func (cfg *config) crawlPage(rawCurrentURL string) {
	// Envia um sinal para o canal de controle de concorrência
	cfg.concurrencyControl <- struct{}{}

	// Quando a função retorna, recebe o sinal enviado antes
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	if cfg.pagesLen() >= cfg.maxPages {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("erro - crawlPage: não foi possível fazer parse da URL '%s': %v\n", rawCurrentURL, err)
		return
	}

	if currentURL.Hostname() != cfg.baseURL.Hostname() {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Printf("erro ao normalizar URL: %v\n", err)
		return
	}

	isFirst := cfg.addPageVisit(normalizedURL)
	if !isFirst {
		return
	}

	// Aplica rate limiting antes de fazer a requisição
	if cfg.rateLimiter != nil {
		if err := cfg.rateLimiter.Wait(cfg.ctx); err != nil {
			log.Printf("erro no rate limiter: %v\n", err)
			return
		}
	}

	// Adiciona delay entre requisições para ser mais "gentil"
	if cfg.requestDelay > 0 {
		time.Sleep(cfg.requestDelay)
	}

	// Atualiza o cliente HTTP se houver rotação de proxies
	if cfg.proxyRotator != nil && cfg.proxyRotator.Count() > 0 {
		Client.Transport = cfg.proxyRotator.GetTransport()
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		log.Printf("erro ao buscar HTML de %s: %v\n", rawCurrentURL, err)
		return
	}

	pageData := extractPageData(html, rawCurrentURL)
	cfg.setPageData(normalizedURL, pageData)

	log.Printf("✓ Crawled: %s (encontrados %d links)\n", rawCurrentURL, len(pageData.OutgoingLinks))

	for _, URL := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(URL)
	}
}
