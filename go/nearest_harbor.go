package tabuamare

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
)

// GetNearestHarbor retorna o porto mais próximo de uma coordenada geográfica
func (c *Client) GetNearestHarbor(ctx context.Context, lat, lng float64) (*NearestHarbor, error) {
	if math.IsNaN(lat) || math.IsInf(lat, 0) {
		return nil, &ValidationError{Field: "lat", Message: "latitude must be a valid number"}
	}

	if math.IsNaN(lng) || math.IsInf(lng, 0) {
		return nil, &ValidationError{Field: "lng", Message: "longitude must be a valid number"}
	}

	if lat < -90 || lat > 90 {
		return nil, &ValidationError{Field: "lat", Message: "latitude must be between -90 and 90 degrees"}
	}

	if lng < -180 || lng > 180 {
		return nil, &ValidationError{Field: "lng", Message: "longitude must be between -180 and 180 degrees"}
	}

	latLng := fmt.Sprintf("%.6f,%.6f", lat, lng)
	path := fmt.Sprintf("/nearest-harbor-independent-state/%s", latLng)

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
