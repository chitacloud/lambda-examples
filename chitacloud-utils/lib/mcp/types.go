// Package mcp provides utilities for creating Model Context Protocol (MCP) servers
package mcp

import "github.com/fredyk/westack-go/v2/lambdas"

// MCPRequest represents a standard MCP protocol request
type MCPRequest struct {
	lambdas.LambdaRequest
	JSONRPC string         `json:"jsonrpc"`
	ID      int            `json:"id"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params"`
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
	Type       string                       `json:"type"`
	Properties map[string]ParameterProperty `json:"properties"`
	Required   []string                     `json:"required"`
}

// ParameterProperty describes a parameter property
type ParameterProperty struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}
