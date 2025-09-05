// Package mcp provides utilities for creating Model Context Protocol (MCP) servers
package mcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

// Server represents an MCP protocol server
type Server struct {
	Name           string
	Version        string
	Description    string
	Tools          []ToolDescription
	DefaultHandler func(r *http.Request, params map[string]any) (any, error)
	Debug          bool
}

// NewServer creates a new MCP server with the given parameters
func NewServer(name, version, description string) *Server {
	return &Server{
		Name:        name,
		Version:     version,
		Description: description,
		Tools:       []ToolDescription{},
		Debug:       false,
	}
}

func (s *Server) SetDebug(debug bool) {
	s.Debug = debug
}

// RegisterTool adds a tool to the server's available tools
func (s *Server) RegisterTool(tool ToolDescription) {
	s.Tools = append(s.Tools, tool)
}

func (s *Server) SetDefaultHandler(handler func(r *http.Request, params map[string]any) (any, error)) {
	s.DefaultHandler = handler
}

func (s *Server) Handle(r *http.Request, w http.ResponseWriter, req MCPRequest) (io.ReadCloser, error) {

	if s.Debug {
		fmt.Println("MCP request:", req.JSONRPC, req.ID, req.Method)
		fmt.Println("Request headers:", r.Header)
		// DEBUG, print body
		fmt.Println("[DEBUG] Request body:", string(req.LambdaRequest.Payload))

	}

	mcpInfo, err := InitHttp(r, w, req)
	if err != nil {
		return nil, err
	} else if mcpInfo.IsPreflight {
		return nil, nil
	}

	// Prepare the response based on path
	var responseData any
	var tool *ToolDescription

	// Handle different MCP protocol paths
	switch mcpInfo.Method {
	case "initialize":
		// Initialize request - return server capabilities
		responseData = s.HandleInitialize()
		if s.Debug {
			fmt.Println("Sending initialize response")
		}

	case "tools/list":
		// List tools request
		responseData = s.HandleTools()
		if s.Debug {
			fmt.Println("Sending tools list response")
		}

	case "tools/call":
		toolName := req.Params.Name

		if tool = s.FindTool(toolName); tool != nil {
			responseData, err = tool.Handler(r, req.Params.Arguments)
			if err != nil {
				fmt.Printf("Error calling tool %s: %s\n", toolName, err.Error())
			}		} else {
			fmt.Printf("Tool %s not found\n", toolName)
			err = errors.New("tool not found")
		}
	default:
		if s.DefaultHandler == nil {
			fmt.Printf("[DEBUG] Default handler not set\n")
			responseData = map[string]any{"status": "OK"}
		} else {
			responseData, err = s.DefaultHandler(r, req.Params.Arguments)
		}

		if s.Debug {
			fmt.Println("Sending default path response")
		}
	}

	return Response(mcpInfo, responseData, err, tool)
}

func wrapToValidToolCallResponse(entry any) (map[string]any, error) {
	unstructuredBytes, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}

	wrappedEntry := map[string]any{
		"content": []map[string]any{
			{"type": "text", "text": string(unstructuredBytes)},
		},
		"structuredContent": entry,
	}
	return wrappedEntry, nil
}

func (s *Server) FindTool(name string) *ToolDescription {
	for _, tool := range s.Tools {
		if tool.Name == name {
			return &tool
		}
	}
	return nil
}

// SetCORSHeaders sets standard CORS headers to allow MCP Inspector to connect
func SetCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// SetStreamableHTTPHeaders sets the headers for MCP Streamable HTTP transport
func SetStreamableHTTPHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/mcp-stream+json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}

// HandleInitialize creates the initialize response data
func (s *Server) HandleInitialize() map[string]any {
	return map[string]any{
		"protocolVersion": "2025-03-26",
		"capabilities": map[string]any{
			"tools": map[string]any{
				"listChanged": true,
			},
		},
		"serverInfo": map[string]any{
			"name":        s.Name,
			"version":     s.Version,
			"description": s.Description,
		},
	}
}

// HandleTools creates the tools list response data
func (s *Server) HandleTools() map[string]interface{} {
	return map[string]interface{}{
		"tools": s.Tools,
	}
}

// GetMethodFromPath extracts the method name from the request path
func GetMethodFromPath(path string) string {
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) > 0 {
		return pathParts[len(pathParts)-1]
	}
	return ""
}

type MCPInfo struct {
	Method      string
	RequestID   int
	IsPreflight bool
	StreamID    string
}

func InitHttp(r *http.Request, w http.ResponseWriter, req MCPRequest) (MCPInfo, error) {

	// Set CORS headers to allow MCP Inspector to connect
	SetCORSHeaders(w)

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return MCPInfo{IsPreflight: true}, nil
	}

	// Set Streamable HTTP headers
	SetStreamableHTTPHeaders(w)

	var err error

	// Define default request ID for JSON-RPC 2.0
	requestID := req.ID

	paramsB, err := json.Marshal(req.Params)
	if err != nil {
		return MCPInfo{RequestID: requestID}, err
	}
	fmt.Printf("Params: %s\n", string(paramsB))

	// Use request method if available, otherwise use path-derived method
	method := req.Method
	if method == "" {
		method = "response"
	}

	return MCPInfo{Method: method, RequestID: req.ID, IsPreflight: false}, nil

}

func Response(mcpInfo MCPInfo, responseData any, err error, tool *ToolDescription) (io.ReadCloser, error) {
	// Use reflection to check if responseData is a slice
	val := reflect.ValueOf(responseData)

	// Format as newline-delimited JSON
	var buffer strings.Builder

	if mcpInfo.Method == "tools/call" && val.Kind() == reflect.Slice {
		// If no streamId, generate a new one
		if mcpInfo.StreamID == "" {
			mcpInfo.StreamID = uuid.New().String()
		}

		if slice, ok := responseData.([]map[string]any); ok {
			// Handle slice streaming by sending each item as a separate JSON-RPC message.
			for i, item := range slice {
				jsonRpcResponse, err := FormatMCPServerResponse(mcpInfo.RequestID, "tools/stream", mcpInfo.StreamID, item, nil)
				if err != nil {
					return nil, fmt.Errorf("failed to format stream data for element %d: %w", i, err)
				}
				buffer.WriteString(string(jsonRpcResponse))
				buffer.WriteString("\n") // Use newline delimiter
			}
		} else {
			// Handle other slice types if necessary, or return an error.
			return nil, fmt.Errorf("unsupported slice type for streaming: %T", responseData)
		}
	} else {
		// Handle as a single, non-streaming response
		responseBody, err := FormatMCPServerResponse(mcpInfo.RequestID, mcpInfo.Method, mcpInfo.StreamID, responseData, err)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}
		buffer.WriteString(string(responseBody))
		buffer.WriteString("\n")
	}

	return io.NopCloser(strings.NewReader(buffer.String())), nil
}
