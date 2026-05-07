package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// ProxyRotator gerencia a rotação de proxies para evitar bloqueios
type ProxyRotator struct {
	proxies     []*url.URL
	currentIdx  int
	mu          sync.Mutex
	healthCheck bool
}

// NewProxyRotator cria um novo rotacionador de proxies
func NewProxyRotator(proxyURLs []string, healthCheck bool) (*ProxyRotator, error) {
	proxies := make([]*url.URL, 0, len(proxyURLs))

	for _, proxyURL := range proxyURLs {
		parsed, err := url.Parse(proxyURL)
		if err != nil {
			log.Printf("Aviso: proxy inválido '%s': %v", proxyURL, err)
			continue
		}
		proxies = append(proxies, parsed)
	}

	if len(proxies) == 0 {
		return nil, fmt.Errorf("nenhum proxy válido fornecido")
	}

	pr := &ProxyRotator{
		proxies:     proxies,
		currentIdx:  0,
		healthCheck: healthCheck,
	}

	if healthCheck {
		pr.removeUnhealthyProxies()
	}

	return pr, nil
}

// GetNextProxy retorna o próximo proxy na rotação
func (pr *ProxyRotator) GetNextProxy() *url.URL {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	if len(pr.proxies) == 0 {
		return nil
	}

	proxy := pr.proxies[pr.currentIdx]
	pr.currentIdx = (pr.currentIdx + 1) % len(pr.proxies)

	return proxy
}

// GetTransport retorna um http.Transport configurado com o próximo proxy
func (pr *ProxyRotator) GetTransport() *http.Transport {
	proxy := pr.GetNextProxy()

	if proxy == nil {
		return &http.Transport{}
	}

	return &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
}

// removeUnhealthyProxies remove proxies que não estão funcionando
func (pr *ProxyRotator) removeUnhealthyProxies() {
	healthyProxies := make([]*url.URL, 0)

	for _, proxy := range pr.proxies {
		if pr.isProxyHealthy(proxy) {
			healthyProxies = append(healthyProxies, proxy)
		} else {
			log.Printf("Proxy não saudável removido: %s", proxy.String())
		}
	}

	pr.proxies = healthyProxies
}

// isProxyHealthy verifica se um proxy está funcionando
func (pr *ProxyRotator) isProxyHealthy(proxy *url.URL) bool {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
		Timeout: 10 * time.Second,
	}

	// Tenta fazer uma requisição simples para verificar o proxy
	resp, err := client.Get("https://www.google.com")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// Count retorna o número de proxies disponíveis
func (pr *ProxyRotator) Count() int {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	return len(pr.proxies)
}
