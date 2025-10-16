//go:build integration
// +build integration

package tabuamare

import (
	"context"
	"testing"
	"time"
)

// Para executar testes de integração: go test -v -tags=integration

func TestIntegration_GetStates(t *testing.T) {
	client := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	states, err := client.GetStates(ctx)
	if err != nil {
		t.Fatalf("failed to get states: %v", err)
	}

	if len(states) == 0 {
		t.Fatal("expected at least one state")
	}

	t.Logf("Found %d states: %v", len(states), states)
}

func TestIntegration_GetHarborNames(t *testing.T) {
	client := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	harbors, err := client.GetHarborNames(ctx, "sc")
	if err != nil {
		t.Fatalf("failed to get harbor names: %v", err)
	}

	if len(harbors) == 0 {
		t.Fatal("expected at least one harbor")
	}

	t.Logf("Found %d harbors in SC", len(harbors))
	for _, h := range harbors {
		t.Logf("  - ID: %d, Name: %s", h.ID, h.HarborName)
	}
}

func TestIntegration_GetHarbor(t *testing.T) {
	client := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	harbor, err := client.GetHarbor(ctx, 1)
	if err != nil {
		t.Fatalf("failed to get harbor: %v", err)
	}

	if harbor.ID != 1 {
		t.Errorf("expected harbor ID 1, got %d", harbor.ID)
	}

	t.Logf("Harbor: %s", harbor.HarborName)
	t.Logf("State: %s", harbor.State)
	t.Logf("Mean Level: %.2f m", harbor.MeanLevel)
}

func TestIntegration_GetTideTable(t *testing.T) {
	client := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tides, err := client.GetTideTable(ctx, 1, 1, []int{1, 2, 3})
	if err != nil {
		t.Fatalf("failed to get tide table: %v", err)
	}

	if len(tides) == 0 {
		t.Fatal("expected at least one tide table")
	}

	t.Logf("Got tide table for %s", tides[0].HarborName)
	for _, month := range tides[0].Months {
		t.Logf("Month: %s", month.MonthName)
		for _, day := range month.Days {
			t.Logf("  Day %d (%s): %d tide records", day.Day, day.WeekdayName, len(day.Hours))
		}
	}
}
