package mcp_hour

import (
	"fmt"
	"io"
	"net/http"

	chitamcputils "github.com/chitacloud/lambda-examples/chitacloud-utils/chita-mcp-utils"
	"github.com/chitacloud/lambda-examples/chitacloud-utils/lib/mcp"
)

// HourResponse represents the response structure
type HourResponse struct {
	Hour        int    `json:"hour"`
	AmPm        string `json:"amPm"`
	Message     string `json:"message"`
	CurrentTime string `json:"currentTime"`
	Error       string `json:"error"`
}

func (hr HourResponse) ToMap() map[string]any {
	return map[string]any{
		"hour":        hr.Hour,
		"amPm":        hr.AmPm,
		"message":     hr.Message,
		"currentTime": hr.CurrentTime,
		"error":       hr.Error,
	}
}

var server *mcp.Server

// Handler is the lambda entry point for Chita Cloud
func Handler(r *http.Request, w http.ResponseWriter, req mcp.MCPRequest) (io.ReadCloser, error) {
	fmt.Println("MCP request:", req.JSONRPC, req.ID, req.Method)
	fmt.Println("Request headers:", r.Header)

	return server.Handle(w, r, req)
}

func init() {
	// Create server with tools
	server = chitamcputils.DefaultServer(chitamcputils.ServerOptions{
		Name:        "HourMCP",
		Version:     "1.0.0",
		Description: "MCP server that provides current timezone",
	})

	registerGetHourTool(server)

	registerGetTimeTool(server)

	registerDefaultHandler(server)
}
