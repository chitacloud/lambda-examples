package mcp_hour

import (
	"fmt"
	"io"
	"net/http"

	chitamcputils "github.com/chitacloud/example-lambdas/chita-utils/chita-mcp-utils"
	"github.com/chitacloud/example-lambdas/chita-utils/lib/mcp"
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

	mcpInfo, err := mcp.InitHttp(r, w, req)
	if err != nil {
		return nil, err
	} else if mcpInfo.IsPreflight {
		return nil, nil
	}

	return server.Handle(w, r, req, mcpInfo)
}

func init() {
	// Create server with tools
	server = chitamcputils.CreateMCPServer()

	registerGetHourTool(server)

	registerGetTimeTool(server)

	registerDefaultHandler(server)
}
