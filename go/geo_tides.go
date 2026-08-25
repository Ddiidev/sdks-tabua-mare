package tabuamare

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// GetGeoTideTable retorna a tábua de marés do porto mais próximo das
// coordenadas informadas dentro do estado especificado
func (c *Client) GetGeoTideTable(ctx context.Context, lat, lng float64, state string, month int, days []int) ([]TideTable, error) {
	if strings.TrimSpace(state) == "" {
		return nil, &ValidationError{Field: "state", Message: "state cannot be empty"}
	}

	if err := validateCoordinates(lat, lng); err != nil {
		return nil, err
	}

	if month < 1 || month > 12 {
		return nil, ErrInvalidMonth
	}

	dayRange, err := NewDayRange(days...)
	if err != nil {
		return nil, err
	}

	state = strings.ToLower(state)
	path := fmt.Sprintf(
		"/geo-tabua-mare/%s/%s/%d/%s",
		formatLatLng(lat, lng),
		state,
		month,
		url.PathEscape(dayRange.String()),
	)

	body, err := c.doRequest(ctx, "GET", path)
	if err != nil {
		return nil, err
	}

	var response TideTableResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return nil, response.Error
	}

	return response.Data, nil
}
