package adapters_test

import (
	"example-lambdas/mcp-hour/adapters"
	"testing"
	"time"
)

func TestSystemClock_GetCurrentHour(t *testing.T) {
	// Create the adapter
	clock := adapters.NewSystemClock("")

	// Get the current hour from the adapter
	hour12, amPm, timeStr, err := clock.GetCurrentHour()
	if err != nil {
		t.Fatalf("GetCurrentHour() returned an error: %v", err)
	}

	// Get the current time directly for verification
	now := time.Now().UTC()

	// Validate hour format (1-12)
	if hour12 < 1 || hour12 > 12 {
		t.Errorf("Hour should be between 1 and 12, got %d", hour12)
	}

	// Validate AM/PM format
	hour24 := now.Hour()
	expectedAmPm := "AM"
	if hour24 >= 12 {
		expectedAmPm = "PM"
	}
	if amPm != expectedAmPm {
		t.Errorf("Expected AM/PM to be %s, got %s", expectedAmPm, amPm)
	}

	// Convert the expected hour to 12-hour format for comparison
	expectedHour := hour24 % 12
	if expectedHour == 0 {
		expectedHour = 12
	}
	if hour12 != expectedHour {
		t.Errorf("Expected hour to be %d, got %d", expectedHour, hour12)
	}

	// Validate that timeStr contains the current year, month and day
	currentYear := now.Format("2006")
	if timeStr[:4] != currentYear {
		t.Errorf("Time string %s should contain the current year %s", timeStr, currentYear)
	}
}
