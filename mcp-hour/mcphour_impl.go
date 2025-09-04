package mcp_hour

import (
	"encoding/json"
	"fmt"

	"github.com/chitacloud/lambda-examples/mcp-hour/adapters"
	"github.com/chitacloud/lambda-examples/mcp-hour/domain"

	"github.com/chitacloud/lambda-examples/chitacloud-utils/lib/mcp"
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

	message := fmt.Sprintf("Current timestamp%s is %s", inTimeZoneDescription, currentTime)

	// currentTime is already formatted as ISO8601 from the adapter

	// Return a response with the hour data
	return HourResponse{
		Hour:        hour,
		AmPm:        amPm,
		Message:     message,
		CurrentTime: currentTime,
	}, nil
}

func registerGetHourTool(server *mcp.Server) {
	server.RegisterTool(mcp.ToolDescription{
		Name:        "get_hour",
		Description: "Get the current timestamp in UTC",
		InputSchema: mcp.Schema{
			Type:       "object",
			Properties: map[string]mcp.ParameterProperty{},
			Required:   []string{},
		},
		OutputSchema: mcp.Schema{
			Type: "object",
			Properties: map[string]mcp.ParameterProperty{
				"hour": {
					Type:        "integer",
					Description: "Current hour in 12-hour format",
				},
				"amPm": {
					Type:        "string",
					Description: "AM or PM indicator",
				},
				"message": {
					Type:        "string",
					Description: "Message containing the current hour and AM/PM indicator",
				},
				"currentTime": {
					Type:        "string",
					Description: "Current time in ISO format",
				},
			},
			Required: []string{"hour", "amPm", "message", "currentTime"},
		},
		Handler: func(params map[string]any) (map[string]any, error) {
			return getFormattedHourInfo(params)
		},
	})
}

func registerGetTimeTool(server *mcp.Server) {
	server.RegisterTool(mcp.ToolDescription{
		Name:        "get_time",
		Description: "Get the current timestamp in the specified timezone",
		InputSchema: mcp.Schema{
			Type: "object",
			Properties: map[string]mcp.ParameterProperty{
				"timezone": {
					Type:        "string",
					Description: "Optional timezone (defaults to system timezone)",
				},
			},
			Required: []string{},
		},
		OutputSchema: mcp.Schema{
			Type: "object",
			Properties: map[string]mcp.ParameterProperty{
				"hour": {
					Type:        "integer",
					Description: "Current hour in 12-hour format",
				},
				"amPm": {
					Type:        "string",
					Description: "AM or PM indicator",
				},
				"message": {
					Type:        "string",
					Description: "Message containing the current hour and AM/PM indicator",
				},
				"currentTime": {
					Type:        "string",
					Description: "Current time in ISO format",
				},
			},
			Required: []string{"hour", "amPm", "message", "currentTime"},
		},
		Handler: func(params map[string]any) (map[string]any, error) {
			return getFormattedHourInfo(params)
		},
	})
}

func registerDefaultHandler(server *mcp.Server) {
	server.SetDefaultHandler(func(params map[string]any) (map[string]any, error) {
		return getFormattedHourInfo(params)
	})
}

func getFormattedHourInfo(params map[string]any) (map[string]any, error) {

	timezone := ""
	if tz, ok := params["timezone"]; ok {
		timezone = tz.(string)
	}

	hourInfo, err := getHourInfo(timezone)
	if err != nil {
		return nil, err
	}

	unstructuredBytes, err := json.Marshal(hourInfo)
	if err != nil {
		return nil, err
	}

	responseData := map[string]any{
		"content": []map[string]any{
			{
				"type": "text",
				"text": string(unstructuredBytes),
			},
		},
		"structuredContent": hourInfo,
	}
	fmt.Println("Sending get_hour response:", hourInfo)
	return responseData, nil
}
