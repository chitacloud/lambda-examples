package domain

// ClockPort defines the interface for getting time information
// This is a secondary/driven port that will be implemented by adapters
type ClockPort interface {
	// GetCurrentHour returns the current hour (1-12), AM/PM indicator, and full time string
	GetCurrentHour() (int, string, string, error)
	GetDayOfWeek() string
}
