package tabuamare

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetUsage retorna o consumo atual de rate-limit da api_key configurada.
// Requer autenticação: configure o cliente com WithAPIKey.
// Este endpoint não consome a cota mensal.
func (c *Client) GetUsage(ctx context.Context) (*UsageInfo, error) {
	if c.apiKey == "" {
		return nil, &ValidationError{Field: "apiKey", Message: "usage requires an API key: use WithAPIKey"}
	}

	body, err := c.doRequest(ctx, "GET", "/usage")
	if err != nil {
		return nil, err
	}

	var response UsageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return nil, response.Error
	}

	if len(response.Data) == 0 {
		return nil, ErrEmptyResponse
	}

	return &response.Data[0], nil
}
