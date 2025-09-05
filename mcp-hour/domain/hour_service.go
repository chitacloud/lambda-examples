package domain

// HourService provides functionality for getting hour information
// This is the core domain service that uses the ClockPort
type HourService struct {
	clock ClockPort
}

// NewHourService creates a new instance of HourService
func NewHourService(clock ClockPort) *HourService {
	return &HourService{
		clock: clock,
	}
}

// GetHourInfo returns the current hour information
func (s *HourService) GetHourInfo() (int, string, string, error) {
	return s.clock.GetCurrentHour()
}

func (s *HourService) GetDayOfWeek() string {
	return s.clock.GetDayOfWeek()
}
