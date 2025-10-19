package logger

import (
	"fmt"
	"strings"
)

type LogRequest struct {
	Level     string `json:"l"`
	Message   string `json:"m"`
	Timestamp string `json:"t"`
	Service   string `json:"s"`
	RequestID string `json:"rid"`
}

func Log(request LogRequest) {
	escapedMsg := strings.ReplaceAll(request.Message, "\"", "\\\"")
	fmt.Printf("t=%s level=%s service=%s rid=%s msg=\"%s\"\n", 
		request.Timestamp, request.Level, request.Service, request.RequestID, escapedMsg)
}
