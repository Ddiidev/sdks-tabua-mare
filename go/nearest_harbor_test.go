package tabuamare

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetNearestHarbor_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/nearest-harbor-independent-state/-23.550520,-46.633308" {
			t.Errorf("expected path to be /nearest-harbor-independent-state/-23.550520,-46.633308, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := NearestHarborResponse{
			Data: []NearestHarbor{
				{
					Harbor: Harbor{
						ID:         1,
						HarborName: "PORTO DE SANTOS",
						State:      "sp",
						Timezone:   "UTC -03.0",
						Card:       "123",
						GeoLocation: []GeoLocation{
							{
								Lat:          "-23.950520",
								Lng:          "-46.333308",
								DecimalLat:   "23° 57' S",
								DecimalLng:   "46° 20' W",
								LatDirection: "s",
								LngDirection: "w",
							},
						},
						MeanLevel: 1.2,
					},
					Distance: 45.5,
				},
			},
			Total: 1,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	harbor, err := client.GetNearestHarbor(context.Background(), -23.550520, -46.633308)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if harbor == nil {
		t.Fatal("expected harbor to be non-nil")
	}

	if harbor.HarborName != "PORTO DE SANTOS" {
		t.Errorf("expected harbor name to be 'PORTO DE SANTOS', got %s", harbor.HarborName)
	}

	if harbor.Distance != 45.5 {
		t.Errorf("expected distance to be 45.5, got %f", harbor.Distance)
	}
}

func TestGetNearestHarbor_InvalidLatitude(t *testing.T) {
	client := NewClient()

	testCases := []struct {
		name string
		lat  float64
		lng  float64
	}{
		{"latitude too high", 91, -46.633308},
		{"latitude too low", -91, -46.633308},
		{"invalid latitude", math.NaN(), -46.633308},
		{"infinite latitude", math.Inf(1), -46.633308},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.GetNearestHarbor(context.Background(), tc.lat, tc.lng)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Errorf("expected ValidationError, got %T", err)
			}
		})
	}
}

func TestGetNearestHarbor_InvalidLongitude(t *testing.T) {
	client := NewClient()

	testCases := []struct {
		name string
		lat  float64
		lng  float64
	}{
		{"longitude too high", -23.550520, 181},
		{"longitude too low", -23.550520, -181},
		{"invalid longitude", -23.550520, math.NaN()},
		{"infinite longitude", -23.550520, math.Inf(1)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.GetNearestHarbor(context.Background(), tc.lat, tc.lng)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Errorf("expected ValidationError, got %T", err)
			}
		})
	}
}

func TestGetNearestHarbor_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := NearestHarborResponse{
			Data:  []NearestHarbor{},
			Total: 0,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	_, err := client.GetNearestHarbor(context.Background(), -23.550520, -46.633308)
	if err != ErrEmptyResponse {
		t.Errorf("expected ErrEmptyResponse, got %v", err)
	}
}
