// Package mcp provides utilities for creating Model Context Protocol (MCP) servers
package mcp

import "github.com/fredyk/westack-go/v2/lambdas"

type MCPRequestParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
	Meta      map[string]any `json:"_meta"`
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
	Name         string `json:"name"`
	Description  string `json:"description"`
	InputSchema  Schema `json:"inputSchema"`
	OutputSchema Schema `json:"outputSchema"`

	Handler func(params map[string]any) (map[string]any, error) `json:"-"`
}

// Schema describes the parameters for a tool
type Schema struct {
	Type                 string         `json:"type"`
	Properties           map[string]any `json:"properties,omitempty"`
	AdditionalProperties bool           `json:"additionalProperties,omitempty"`
	// for arrays:
	Items    *Schema  `json:"items,omitempty"`
	Required []string `json:"required"`
}
