package tabuamare

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHarbors_BuildsBracketedList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/harbors/[pb01,pe01]" {
			t.Errorf("expected path to be /harbors/[pb01,pe01], got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(HarborsResponse{
			Data: []Harbor{
				{ID: "pb01", HarborName: "PORTO DE CABEDELO", State: "pb"},
				{ID: "pe01", HarborName: "FERNANDO DE NORONHA", State: "pe"},
			},
			Total: 2,
		})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	harbors, err := client.GetHarbors(context.Background(), "pb01", "pe01")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(harbors) != 2 {
		t.Fatalf("expected 2 harbors, got %d", len(harbors))
	}

	if harbors[0].ID != "pb01" || harbors[1].ID != "pe01" {
		t.Errorf("unexpected harbor IDs: %s, %s", harbors[0].ID, harbors[1].ID)
	}
}

func TestGetHarbors_EmptyID(t *testing.T) {
	client := NewClient()

	testCases := [][]string{
		{},
		{""},
		{"pb01", "  "},
	}

	for _, ids := range testCases {
		if _, err := client.GetHarbors(context.Background(), ids...); err == nil {
			t.Errorf("expected error for ids %v, got nil", ids)
		}
	}
}

func TestGetTideTable_Path(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/tabua-mare/pb01/1/%5B1%2C2%5D"
		if r.URL.EscapedPath() != expected {
			t.Errorf("expected path %s, got %s", expected, r.URL.EscapedPath())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(TideTableResponse{
			Data:  []TideTable{{Year: 2026, HarborName: "PORTO DE CABEDELO", State: "pb"}},
			Total: 1,
		})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	tables, err := client.GetTideTable(context.Background(), "pb01", 1, []int{1, 2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tables) != 1 || tables[0].HarborName != "PORTO DE CABEDELO" {
		t.Fatalf("unexpected tables: %+v", tables)
	}
}

func TestGetTideTable_Validation(t *testing.T) {
	client := NewClient()

	if _, err := client.GetTideTable(context.Background(), "", 1, []int{1}); err == nil {
		t.Error("expected error for empty harbor ID")
	}

	if _, err := client.GetTideTable(context.Background(), "pb01", 13, []int{1}); err != ErrInvalidMonth {
		t.Errorf("expected ErrInvalidMonth, got %v", err)
	}
}

func TestGetGeoTideTable_Path(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/geo-tabua-mare/[-7.115090,-34.864000]/pb/1/%5B1%5D"
		if r.URL.EscapedPath() != expected {
			t.Errorf("expected path %s, got %s", expected, r.URL.EscapedPath())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(TideTableResponse{
			Data:  []TideTable{{Year: 2026, HarborName: "PORTO DE CABEDELO", State: "pb"}},
			Total: 1,
		})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	tables, err := client.GetGeoTideTable(context.Background(), -7.11509, -34.864, "PB", 1, []int{1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tables) != 1 || tables[0].State != "pb" {
		t.Fatalf("unexpected tables: %+v", tables)
	}
}

func TestGetGeoTideTable_Validation(t *testing.T) {
	client := NewClient()

	if _, err := client.GetGeoTideTable(context.Background(), -7.11, -34.86, "", 1, []int{1}); err == nil {
		t.Error("expected error for empty state")
	}

	if _, err := client.GetGeoTideTable(context.Background(), -7.11, -34.86, "pb", 0, []int{1}); err != ErrInvalidMonth {
		t.Errorf("expected ErrInvalidMonth, got %v", err)
	}

	_, err := client.GetGeoTideTable(context.Background(), 91, -34.86, "pb", 1, []int{1})
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Errorf("expected ValidationError for invalid lat, got %T", err)
	}
}

func TestGetUsage_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/usage" {
			t.Errorf("expected path /usage, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer tm_live_test" {
			t.Error("expected Authorization header")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(UsageResponse{
			Data: []UsageInfo{{
				Plan:             "plan5",
				LimitRPM:         "512",
				UsedRPM:          "3",
				RemainingRPM:     "509",
				LimitMonthly:     "256000",
				UsedMonthly:      "842",
				RemainingMonthly: "-1",
			}},
			Total: 1,
		})
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL), WithAPIKey("tm_live_test"))
	usage, err := client.GetUsage(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if usage.Plan != "plan5" {
		t.Errorf("expected plan plan5, got %s", usage.Plan)
	}

	if usage.RemainingMonthly != "-1" {
		t.Errorf("expected remaining_monthly '-1' (unlimited), got %s", usage.RemainingMonthly)
	}
}

func TestGetUsage_RequiresAPIKey(t *testing.T) {
	client := NewClient()

	_, err := client.GetUsage(context.Background())
	if err == nil {
		t.Fatal("expected error without API key, got nil")
	}

	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}
