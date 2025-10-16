package tabuamare

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// GetHarborNames retorna a lista de portos de um estado específico
func (c *Client) GetHarborNames(ctx context.Context, state string) ([]HarborName, error) {
	if state == "" {
		return nil, &ValidationError{Field: "state", Message: "state cannot be empty"}
	}

	state = strings.ToLower(state)
	path := fmt.Sprintf("/harbor_names/%s", state)

	body, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var response HarborNamesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return nil, response.Error
	}

	return response.Data, nil
}

// GetHarbors retorna informações detalhadas de um ou mais portos por IDs
func (c *Client) GetHarbors(ctx context.Context, ids ...int) ([]Harbor, error) {
	if len(ids) == 0 {
		return nil, &ValidationError{Field: "ids", Message: "at least one harbor ID is required"}
	}

	for _, id := range ids {
		if id <= 0 {
			return nil, &ValidationError{Field: "ids", Message: "harbor IDs must be positive integers"}
		}
	}

	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = strconv.Itoa(id)
	}

	path := fmt.Sprintf("/harbors/%s", strings.Join(idsStr, ","))

	body, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var response HarborsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return nil, response.Error
	}

	return response.Data, nil
}

// GetHarbor retorna informações detalhadas de um porto específico
func (c *Client) GetHarbor(ctx context.Context, id int) (*Harbor, error) {
	harbors, err := c.GetHarbors(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(harbors) == 0 {
		return nil, ErrEmptyResponse
	}

	return &harbors[0], nil
}
