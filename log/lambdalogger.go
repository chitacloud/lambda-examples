package logger

import "fmt"

type LogRequest struct {
	Level     string `json:"l"`
	Message   string `json:"m"`
	Timestamp string `json:"t"`
	Service   string `json:"s"`
	RequestID string `json:"rid"`
}

func Log(request LogRequest) {
	fmt.Printf("t=%s level=%s service=%s rid=%s msg=\"%s\"\n", 
		request.Timestamp, request.Level, request.Service, request.RequestID, request.Message)
}
