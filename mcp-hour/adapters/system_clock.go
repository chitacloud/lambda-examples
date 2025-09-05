package adapters

import (
	"fmt"
	"time"
)

// SystemClock is an adapter that implements the ClockPort interface
type SystemClock struct {
	timezone string
}

// NewSystemClock creates a new SystemClock adapter
func NewSystemClock(timezone string) *SystemClock {
	return &SystemClock{timezone: timezone}
}

// GetCurrentHour returns the current hour in 12-hour format, AM/PM designation, and full time string
func (c *SystemClock) GetCurrentHour() (int, string, string, error) {
	now := time.Now()

	fmt.Printf("Using timezone: %s\n", c.timezone)

	now, i, s, s1, err := c.getTimeNow(now)
	if err != nil {
		return i, s, s1, err
	}

	// Format the full time string
	var currentTime string
	if c.timezone != "" {
		currentTime = now.Format(time.RFC3339)
	} else {
		currentTime = now.UTC().Format(time.RFC3339)
	}

	// Get hour in 24-hour format
	hour24 := now.Hour()

	// Convert to 12-hour format
	hour12 := hour24 % 12
	if hour12 == 0 {
		hour12 = 12
	}

	// Determine AM/PM
	amPm := "AM"
	if hour24 >= 12 {
		amPm = "PM"
	}

	return hour12, amPm, currentTime, nil
}

func (c *SystemClock) getTimeNow(now time.Time) (time.Time, int, string, string, error) {
	if c.timezone != "" {
		location, err := time.LoadLocation(c.timezone)
		if err != nil {
			return time.Time{}, 0, "", "", fmt.Errorf("failed to load location for tz=%s: %w", c.timezone, err)
		}
		now = now.In(location)
	} else {
		now = now.UTC()
	}
	return now, 0, "", "", nil
}

func (c *SystemClock) GetDayOfWeek() string {

	now, _, _, _, _ := c.getTimeNow(time.Now())
	return now.Weekday().String()
}
