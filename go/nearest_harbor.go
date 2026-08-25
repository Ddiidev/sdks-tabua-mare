package tabuamare

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// GetNearestHarbor retorna o porto mais próximo de uma coordenada geográfica,
// sem restringir por estado
func (c *Client) GetNearestHarbor(ctx context.Context, lat, lng float64) (*Harbor, error) {
	if err := validateCoordinates(lat, lng); err != nil {
		return nil, err
	}

	latLng := formatLatLng(lat, lng)
	path := fmt.Sprintf("/nearest-harbor-independent-state/%s", latLng)

	harbor, err := c.requestNearestHarbor(ctx, path)
	if err != nil {
		return nil, err
	}

	return harbor, nil
}

// GetNearestHarborByState retorna o porto mais próximo de uma coordenada
// geográfica dentro do estado especificado
func (c *Client) GetNearestHarborByState(ctx context.Context, state string, lat, lng float64) (*Harbor, error) {
	if strings.TrimSpace(state) == "" {
		return nil, &ValidationError{Field: "state", Message: "state cannot be empty"}
	}

	if err := validateCoordinates(lat, lng); err != nil {
		return nil, err
	}

	state = strings.ToLower(state)
	latLng := formatLatLng(lat, lng)
	path := fmt.Sprintf("/nearested-harbor/%s/%s", state, latLng)

	harbor, err := c.requestNearestHarbor(ctx, path)
	if err != nil {
		return nil, err
	}

	return harbor, nil
}

// requestNearestHarbor executa a requisição e extrai o porto da resposta
func (c *Client) requestNearestHarbor(ctx context.Context, path string) (*Harbor, error) {
	body, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var response NearestHarborResponse
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

// validateCoordinates valida latitude e longitude
func validateCoordinates(lat, lng float64) error {
	if math.IsNaN(lat) || math.IsInf(lat, 0) {
		return &ValidationError{Field: "lat", Message: "latitude must be a valid number"}
	}

	if math.IsNaN(lng) || math.IsInf(lng, 0) {
		return &ValidationError{Field: "lng", Message: "longitude must be a valid number"}
	}

	if lat < -90 || lat > 90 {
		return &ValidationError{Field: "lat", Message: "latitude must be between -90 and 90 degrees"}
	}

	if lng < -180 || lng > 180 {
		return &ValidationError{Field: "lng", Message: "longitude must be between -180 and 180 degrees"}
	}

	return nil
}

// formatLatLng formata as coordenadas no formato [lat,lng] esperado pela API
func formatLatLng(lat, lng float64) string {
	return fmt.Sprintf("[%.6f,%.6f]", lat, lng)
}
