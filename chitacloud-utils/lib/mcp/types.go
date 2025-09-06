// Package mcp provides utilities for creating Model Context Protocol (MCP) servers
package mcp

import (
	"net/http"

	"github.com/fredyk/westack-go/lambdas"
	"github.com/getkin/kin-openapi/openapi3"
)

type MCPRequestParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
	Meta      map[string]any `json:"_meta"`
	StreamID  string         `json:"streamId,omitempty"`
}

// MCPRequest represents a standard MCP protocol request
type MCPRequest struct {
	lambdas.LambdaRequest
	JSONRPC string           `json:"jsonrpc"`
	ID      int              `json:"id"`
	Method  string           `json:"method"`
	Params  MCPRequestParams `json:"params"`
}

// ToolDescription represents an MCP tool description
type ToolDescription struct {
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	InputSchema  *openapi3.Schema `json:"inputSchema"`
	OutputSchema *openapi3.Schema `json:"outputSchema"`
	Raw          bool             `json:"raw,omitempty"`

	Handler func(r *http.Request, params map[string]any) (any, error) `json:"-"`
}
