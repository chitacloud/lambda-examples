// Package mcp provides utilities for creating Model Context Protocol (MCP) servers
package mcp

import "encoding/json"

// FormatMCPServerResponse formats the response according to JSON-RPC 2.0 / MCP protocol
func FormatMCPServerResponse(id int, method string, content any) ([]byte, error) {
	responseObj := map[string]interface{}{
		"jsonrpc": "2.0",
	}

	// if id != "" {
	responseObj["id"] = id
	// }

	// According to JSON-RPC 2.0, we should use 'result' to contain the response content
	responseObj["result"] = content

	return json.Marshal(responseObj)
}
