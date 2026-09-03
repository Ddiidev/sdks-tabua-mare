package tabuamare

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetStates retorna a lista de todos os estados costeiros brasileiros disponíveis
func (c *Client) GetStates(ctx context.Context) ([]string, error) {
	body, err := c.doRequest(ctx, "GET", "/states")
	if err != nil {
		return nil, err
	}

	var response StatesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return nil, response.Error
	}

	return response.Data, nil
}
