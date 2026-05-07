package main

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter controla a taxa de requisições para evitar bloqueios
type RateLimiter struct {
	limiter *rate.Limiter
}

// NewRateLimiter cria um novo rate limiter
// requestsPerSecond: número de requisições permitidas por segundo
// burst: número máximo de requisições em rajada
func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), burst),
	}
}

// Wait aguarda até que seja permitido fazer uma nova requisição
func (rl *RateLimiter) Wait(ctx context.Context) error {
	return rl.limiter.Wait(ctx)
}

// Allow verifica se uma requisição pode ser feita imediatamente
func (rl *RateLimiter) Allow() bool {
	return rl.limiter.Allow()
}

// DelayBetweenRequests adiciona um delay fixo entre requisições
// Útil para ser mais "gentil" com o servidor alvo
func DelayBetweenRequests(minDelay, maxDelay time.Duration) time.Duration {
	// Implementação simples: retorna o delay mínimo
	// Pode ser expandido para usar delay aleatório
	return minDelay
}
