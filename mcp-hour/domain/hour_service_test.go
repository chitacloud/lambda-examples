package domain_test

import (
	"mcp-hour/domain"
	"testing"
)

// MockClock is a mock implementation of the ClockPort interface
type MockClock struct {
	Hour        int
	AmPm        string
	CurrentTime string
}

func (m *MockClock) GetCurrentHour() (int, string, string) {
	return m.Hour, m.AmPm, m.CurrentTime
}

func TestHourService(t *testing.T) {
	// Create a mock clock and preset the values
	mockClock := &MockClock{
		Hour:        10,
		AmPm:        "AM",
		CurrentTime: "2023-07-15 10:30:00",
	}

	// Create the service with the mock
	hourService := domain.NewHourService(mockClock)

	// Get the hour information
	hour, amPm, currentTime := hourService.GetHourInfo()

	// Verify the returned values
	if hour != 10 {
		t.Errorf("Expected hour to be 10, got %d", hour)
	}

	if amPm != "AM" {
		t.Errorf("Expected AM/PM to be AM, got %s", amPm)
	}

	if currentTime != "2023-07-15 10:30:00" {
		t.Errorf("Expected currentTime to be 2023-07-15 10:30:00, got %s", currentTime)
	}
}
