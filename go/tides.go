package tabuamare

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// DayRange representa um intervalo de dias ou dias específicos
type DayRange struct {
	days []int
}

// NewDayRange cria um novo DayRange a partir de dias específicos
func NewDayRange(days ...int) (*DayRange, error) {
	if len(days) == 0 {
		return nil, &ValidationError{Field: "days", Message: "at least one day is required"}
	}

	for _, day := range days {
		if day < 1 || day > 31 {
			return nil, &ValidationError{Field: "days", Message: "days must be between 1 and 31"}
		}
	}

	return &DayRange{days: days}, nil
}

// NewDayRangeFromInterval cria um DayRange a partir de um intervalo (ex: 1-15)
func NewDayRangeFromInterval(start, end int) (*DayRange, error) {
	if start < 1 || start > 31 {
		return nil, &ValidationError{Field: "start", Message: "start day must be between 1 and 31"}
	}
	if end < 1 || end > 31 {
		return nil, &ValidationError{Field: "end", Message: "end day must be between 1 and 31"}
	}
	if start > end {
		return nil, &ValidationError{Field: "days", Message: "start day must be less than or equal to end day"}
	}

	days := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		days = append(days, i)
	}

	return &DayRange{days: days}, nil
}

// String retorna a representação em string do DayRange no formato esperado pela API
func (dr *DayRange) String() string {
	if len(dr.days) == 0 {
		return "[]"
	}

	parts := make([]string, len(dr.days))
	for i, day := range dr.days {
		parts[i] = strconv.Itoa(day)
	}

	return fmt.Sprintf("[%s]", strings.Join(parts, ","))
}

// GetTideTable retorna a tábua de marés para um porto, mês e dias específicos
func (c *Client) GetTideTable(ctx context.Context, harborID, month int, days []int) ([]TideTable, error) {
	if harborID <= 0 {
		return nil, &ValidationError{Field: "harborID", Message: "harbor ID must be a positive integer"}
	}

	if month < 1 || month > 12 {
		return nil, ErrInvalidMonth
	}

	dayRange, err := NewDayRange(days...)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/tabua-mare/%d/%d/%s", harborID, month, url.PathEscape(dayRange.String()))

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

// GetTideTableForMonth retorna a tábua de marés para um mês inteiro
func (c *Client) GetTideTableForMonth(ctx context.Context, harborID, month int) ([]TideTable, error) {
	dayRange, err := NewDayRangeFromInterval(1, 31)
	if err != nil {
		return nil, err
	}

	if harborID <= 0 {
		return nil, &ValidationError{Field: "harborID", Message: "harbor ID must be a positive integer"}
	}

	if month < 1 || month > 12 {
		return nil, ErrInvalidMonth
	}

	path := fmt.Sprintf("/tabua-mare/%d/%d/%s", harborID, month, url.PathEscape(dayRange.String()))

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
