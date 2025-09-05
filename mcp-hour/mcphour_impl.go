package mcp_hour

import (
	"fmt"
	"net/http"

	"github.com/chitacloud/lambda-examples/chitacloud-utils/lib/mcp"
	"github.com/chitacloud/lambda-examples/mcp-hour/adapters"
	"github.com/chitacloud/lambda-examples/mcp-hour/domain"
	"github.com/getkin/kin-openapi/openapi3"
)

// getHourInfo returns the current hour information using domain services
func getHourInfo(timezone string) (HourResponse, error) {
	// Create the domain service with a system clock adapter
	clockAdapter := adapters.NewSystemClock(timezone)
	hourService := domain.NewHourService(clockAdapter)

	// Get the hour information from domain service
	hour, amPm, currentTime, err := hourService.GetHourInfo()
	if err != nil {
		return HourResponse{}, err
	}

	var inTimeZoneDescription string
	if timezone != "" {
		inTimeZoneDescription = fmt.Sprintf(" in %s", timezone)
	}

	dayOfWeek := hourService.GetDayOfWeek()

	message := fmt.Sprintf("Current timestamp%s is %s. Today is %s.", inTimeZoneDescription, currentTime, dayOfWeek)

	// currentTime is already formatted as ISO8601 from the adapter

	// Return a response with the hour data
	return HourResponse{
		Hour:        hour,
		AmPm:        amPm,
		Message:     message,
		CurrentTime: currentTime,
		DayOfWeek:   dayOfWeek,
	}, nil
}

func registerGetTimeTool(server *mcp.Server) {
	server.RegisterTool(mcp.ToolDescription{
		Name:        "get_time",
		Description: "Get the current timestamp in the specified timezone",
		InputSchema: &openapi3.Schema{
			Type: &openapi3.Types{openapi3.TypeObject},
			Properties: map[string]*openapi3.SchemaRef{
				"timezone": {
					Value: &openapi3.Schema{
						Type:        &openapi3.Types{openapi3.TypeString},
						Description: "Optional timezone (defaults to system timezone)",
					},
				},
			},
		},
		OutputSchema: &openapi3.Schema{
			Type: &openapi3.Types{openapi3.TypeObject},
			Properties: map[string]*openapi3.SchemaRef{
				"hour": {
					Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeInteger}, Description: "Current hour in 12-hour format"},
				},
				"amPm": {
					Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}, Description: "AM or PM indicator"},
				},
				"message": {
					Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}, Description: "Message containing the current hour and AM/PM indicator"},
				},
				"dayOfWeek": {
					Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}, Description: "Day of the week"},
				},
				"currentTime": {
					Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}, Description: "Current time in ISO format"},
				},
			},
			Required: []string{"hour", "amPm", "message", "currentTime", "dayOfWeek"},
		},
		Handler: func(r *http.Request, params map[string]any) (any, error) {
			return getFormattedHourInfo(params)
		},
	})
}

func registerDefaultHandler(server *mcp.Server) {
	server.SetDefaultHandler(func(r *http.Request, params map[string]any) (any, error) {
		return getFormattedHourInfo(params)
	})
}

func getFormattedHourInfo(params map[string]any) (any, error) {

	timezone := ""
	if tz, ok := params["timezone"]; ok {
		timezone = tz.(string)
	}

	hourInfo, err := getHourInfo(timezone)
	if err != nil {
		return nil, err
	}

	fmt.Println("Sending get_hour response:", hourInfo)
	return hourInfo.ToMap(), nil
}
