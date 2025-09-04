package mcp_hour

import (
	"example-lambdas/mcp-hour/adapters"
	"example-lambdas/mcp-hour/domain"
	"fmt"

	"github.com/chitacloud/chita-utils/lib/mcp"
)

// getHourInfo returns the current hour information using domain services
func getHourInfo(timezone string) HourResponse {
	// Create the domain service with a system clock adapter
	clockAdapter := adapters.NewSystemClock(timezone)
	hourService := domain.NewHourService(clockAdapter)

	// Get the hour information from domain service
	hour, amPm, currentTime, err := hourService.GetHourInfo()
	if err != nil {
		return HourResponse{Error: err.Error()}
	}
	message := fmt.Sprintf("Current hour is %s", currentTime)

	// currentTime is already formatted as ISO8601 from the adapter

	// Return a response with the hour data
	return HourResponse{
		Hour:        hour,
		AmPm:        amPm,
		Message:     message,
		CurrentTime: currentTime,
	}
}

func registerGetHourTool(server *mcp.Server) {
	server.RegisterTool(mcp.ToolDescription{
		Name:        "get_hour",
		Description: "Get the current hour",
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
			responseData := getFormattedHourInfo(params)
			return responseData, nil
		},
	})
}

func registerGetTimeTool(server *mcp.Server) {
	server.RegisterTool(mcp.ToolDescription{
		Name:        "get_time",
		Description: "Get the current time",
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
			responseData := getFormattedHourInfo(params)
			return responseData, nil
		},
	})
}

func registerDefaultHandler(server *mcp.Server) {
	server.SetDefaultHandler(func(params map[string]any) (map[string]any, error) {
		return getFormattedHourInfo(params), nil
	})
}

func getFormattedHourInfo(params map[string]any) map[string]any {

	timezone := ""
	if tz, ok := params["timezone"]; ok {
		timezone = tz.(string)
	}

	hourInfo := getHourInfo(timezone)
	if hourInfo.Error != "" {
		return map[string]any{"error": hourInfo.Error}
	}
	responseData := map[string]any{
		"content": []map[string]any{
			{
				"type": "text",
				"text": hourInfo.Message,
			},
		},
		"structuredContent": hourInfo,
	}
	fmt.Println("Sending get_hour response:", hourInfo)
	return responseData
}
