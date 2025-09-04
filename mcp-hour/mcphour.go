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

var server *mcp.Server

// Handler is the lambda entry point for Chita Cloud
func Handler(r *http.Request, w http.ResponseWriter, req mcp.MCPRequest) (io.ReadCloser, error) {
	fmt.Println("MCP request:", req.JSONRPC, req.ID, req.Method)
	fmt.Println("Request headers:", r.Header)

	return server.Handle(w, r, req)
}

func init() {
	// Create server with tools
	server = chitamcputils.CreateMCPServer()

	registerGetHourTool(server)

	registerGetTimeTool(server)

	registerDefaultHandler(server)
}
