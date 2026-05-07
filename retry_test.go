package main

import (
	"errors"
	"testing"
	"time"
)

func TestDefaultRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()

	if config.MaxRetries != 3 {
		t.Errorf("MaxRetries esperado 3, obtido %d", config.MaxRetries)
	}

	if config.InitialBackoff != 1*time.Second {
		t.Errorf("InitialBackoff esperado 1s, obtido %v", config.InitialBackoff)
	}

	if config.MaxBackoff != 30*time.Second {
		t.Errorf("MaxBackoff esperado 30s, obtido %v", config.MaxBackoff)
	}

	if config.Multiplier != 2.0 {
		t.Errorf("Multiplier esperado 2.0, obtido %f", config.Multiplier)
	}
}

func TestRetryWithBackoffSuccess(t *testing.T) {
	config := RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
	}

	attempts := 0
	operation := func() error {
		attempts++
		return nil // Sucesso na primeira tentativa
	}

	err := RetryWithBackoff(config, operation)

	if err != nil {
		t.Errorf("Não esperava erro: %v", err)
	}

	if attempts != 1 {
		t.Errorf("Esperava 1 tentativa, obteve %d", attempts)
	}
}

func TestRetryWithBackoffEventualSuccess(t *testing.T) {
	config := RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
	}

	attempts := 0
	operation := func() error {
		attempts++
		if attempts < 3 {
			return errors.New("falha temporária")
		}
		return nil // Sucesso na terceira tentativa
	}

	err := RetryWithBackoff(config, operation)

	if err != nil {
		t.Errorf("Não esperava erro: %v", err)
	}

	if attempts != 3 {
		t.Errorf("Esperava 3 tentativas, obteve %d", attempts)
	}
}

func TestRetryWithBackoffMaxRetriesExceeded(t *testing.T) {
	config := RetryConfig{
		MaxRetries:     2,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
	}

	attempts := 0
	operation := func() error {
		attempts++
		return errors.New("falha permanente")
	}

	err := RetryWithBackoff(config, operation)

	if err == nil {
		t.Error("Esperava erro após max retries")
	}

	// MaxRetries + 1 tentativa inicial
	expectedAttempts := config.MaxRetries + 1
	if attempts != expectedAttempts {
		t.Errorf("Esperava %d tentativas, obteve %d", expectedAttempts, attempts)
	}
}

func TestCalculateBackoff(t *testing.T) {
	config := RetryConfig{
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     10 * time.Second,
		Multiplier:     2.0,
	}

	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, 1 * time.Second},  // 1 * 2^0 = 1
		{1, 2 * time.Second},  // 1 * 2^1 = 2
		{2, 4 * time.Second},  // 1 * 2^2 = 4
		{3, 8 * time.Second},  // 1 * 2^3 = 8
		{4, 10 * time.Second}, // 1 * 2^4 = 16, mas limitado a 10
	}

	for _, tt := range tests {
		backoff := calculateBackoff(config, tt.attempt)
		if backoff != tt.expected {
			t.Errorf("Tentativa %d: esperava %v, obteve %v", tt.attempt, tt.expected, backoff)
		}
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{errors.New("algum erro"), true},
	}

	for _, tt := range tests {
		result := IsRetryableError(tt.err)
		if result != tt.expected {
			t.Errorf("IsRetryableError(%v) = %v, esperava %v", tt.err, result, tt.expected)
		}
	}
}

func TestRetryWithBackoffTiming(t *testing.T) {
	config := RetryConfig{
		MaxRetries:     2,
		InitialBackoff: 50 * time.Millisecond,
		MaxBackoff:     200 * time.Millisecond,
		Multiplier:     2.0,
	}

	attempts := 0
	operation := func() error {
		attempts++
		if attempts < 3 {
			return errors.New("falha")
		}
		return nil
	}

	start := time.Now()
	err := RetryWithBackoff(config, operation)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Não esperava erro: %v", err)
	}

	// Tempo mínimo esperado: 50ms + 100ms = 150ms
	minExpected := 150 * time.Millisecond
	if elapsed < minExpected {
		t.Errorf("Tempo decorrido %v é menor que o esperado %v", elapsed, minExpected)
	}
}
