package main

import (
	"context"
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(10.0, 5)

	if rl == nil {
		t.Fatal("NewRateLimiter retornou nil")
	}

	if rl.limiter == nil {
		t.Fatal("limiter interno é nil")
	}
}

func TestRateLimiterAllow(t *testing.T) {
	rl := NewRateLimiter(10.0, 5)

	// Primeiras 5 requisições devem ser permitidas (burst)
	for i := 0; i < 5; i++ {
		if !rl.Allow() {
			t.Errorf("Requisição %d deveria ser permitida (burst)", i+1)
		}
	}
}

func TestRateLimiterWait(t *testing.T) {
	rl := NewRateLimiter(100.0, 1) // 100 req/s para teste rápido
	ctx := context.Background()

	start := time.Now()

	// Primeira requisição deve ser imediata
	err := rl.Wait(ctx)
	if err != nil {
		t.Fatalf("Wait falhou: %v", err)
	}

	elapsed := time.Since(start)
	if elapsed > 10*time.Millisecond {
		t.Errorf("Primeira requisição demorou muito: %v", elapsed)
	}
}

func TestRateLimiterWaitWithContext(t *testing.T) {
	rl := NewRateLimiter(1.0, 1) // 1 req/s

	// Consome o burst
	rl.Allow()

	// Cria contexto com timeout curto
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Segunda requisição deve esperar, mas contexto vai cancelar
	err := rl.Wait(ctx)
	if err == nil {
		t.Error("Esperava erro de contexto cancelado")
	}
}

func TestDelayBetweenRequests(t *testing.T) {
	minDelay := 100 * time.Millisecond
	maxDelay := 500 * time.Millisecond

	delay := DelayBetweenRequests(minDelay, maxDelay)

	if delay < minDelay {
		t.Errorf("Delay %v é menor que minDelay %v", delay, minDelay)
	}

	if delay > maxDelay {
		t.Errorf("Delay %v é maior que maxDelay %v", delay, maxDelay)
	}
}
