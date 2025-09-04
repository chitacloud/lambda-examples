// Package mcp provides utilities for creating Model Context Protocol (MCP) servers
package mcp

import "encoding/json"

const (
	ErrUnkown = -32001
)

// JsonRPCError represents a JSON-RPC 2.0 error
type JsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// FormatMCPServerResponse formats the response according to JSON-RPC 2.0 / MCP protocol
func FormatMCPServerResponse(id int, method string, content any, err error) ([]byte, error) {
	responseObj := map[string]any{
		"jsonrpc": "2.0",
	}

	responseObj["id"] = id

	if err != nil {
		responseObj["error"] = JsonRPCError{Code: ErrUnkown, Message: err.Error(), Data: map[string]any{"content": content}}
	} else {
		// According to JSON-RPC 2.0, we should use 'result' to contain the response content
		responseObj["result"] = content
	}

	return json.Marshal(responseObj)
}
