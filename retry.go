package main

import (
	"fmt"
	"log"
	"math"
	"time"
)

// RetryConfig define a configuração para retry automático
type RetryConfig struct {
	MaxRetries     int           // Número máximo de tentativas
	InitialBackoff time.Duration // Tempo inicial de espera
	MaxBackoff     time.Duration // Tempo máximo de espera
	Multiplier     float64       // Multiplicador para backoff exponencial
}

// DefaultRetryConfig retorna uma configuração padrão de retry
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
		Multiplier:     2.0,
	}
}

// RetryWithBackoff executa uma função com retry automático e backoff exponencial
// Útil para lidar com falhas temporárias de rede ou rate limiting do servidor
func RetryWithBackoff(config RetryConfig, operation func() error) error {
	var err error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		err = operation()

		if err == nil {
			return nil
		}

		if attempt == config.MaxRetries {
			return fmt.Errorf("falhou após %d tentativas: %w", config.MaxRetries+1, err)
		}

		// Calcula o tempo de espera com backoff exponencial
		backoff := calculateBackoff(config, attempt)
		log.Printf("Tentativa %d falhou: %v. Aguardando %v antes de tentar novamente...",
			attempt+1, err, backoff)

		time.Sleep(backoff)
	}

	return err
}

// calculateBackoff calcula o tempo de espera usando backoff exponencial
func calculateBackoff(config RetryConfig, attempt int) time.Duration {
	backoff := float64(config.InitialBackoff) * math.Pow(config.Multiplier, float64(attempt))

	if backoff > float64(config.MaxBackoff) {
		return config.MaxBackoff
	}

	return time.Duration(backoff)
}

// IsRetryableError verifica se um erro deve ser retentado
// Pode ser expandido para incluir mais casos específicos
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Adicione aqui lógica para identificar erros que devem ser retentados
	// Por exemplo: timeouts, erros 429 (Too Many Requests), 503 (Service Unavailable)
	return true
}
