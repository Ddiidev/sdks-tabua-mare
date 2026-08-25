package tabuamare

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrInvalidState é retornado quando o estado fornecido é inválido
	ErrInvalidState = errors.New("invalid state code")

	// ErrInvalidHarborID é retornado quando o ID do porto é inválido
	ErrInvalidHarborID = errors.New("invalid harbor ID")

	// ErrInvalidMonth é retornado quando o mês está fora do intervalo 1-12
	ErrInvalidMonth = errors.New("invalid month: must be between 1 and 12")

	// ErrInvalidDays é retornado quando os dias fornecidos são inválidos
	ErrInvalidDays = errors.New("invalid days")

	// ErrEmptyResponse é retornado quando a API retorna uma resposta vazia
	ErrEmptyResponse = errors.New("empty response from API")

	// ErrRateLimitExceeded é retornado quando o limite de requisições é excedido
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	// ErrInvalidCoordinates é retornado quando as coordenadas geográficas são inválidas
	ErrInvalidCoordinates = errors.New("invalid coordinates")
)

// APIError representa um erro retornado pela API
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int    // HTTP status code
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d, code %d): %s", e.Status, e.Code, e.Message)
}

// RateLimitError é retornado quando a API responde 429.
// Implementa errors.Is com ErrRateLimitExceeded e expõe o header Retry-After.
type RateLimitError struct {
	RetryAfter time.Duration // tempo a aguardar antes de tentar novamente
}

func (e *RateLimitError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("rate limit exceeded: retry after %s", e.RetryAfter)
	}
	return ErrRateLimitExceeded.Error()
}

func (e *RateLimitError) Is(target error) bool {
	return target == ErrRateLimitExceeded
}

// IsAPIError verifica se um erro é do tipo APIError
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}

// NetworkError representa um erro de rede
type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error: %v", e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

// ValidationError representa um erro de validação de parâmetros
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}
