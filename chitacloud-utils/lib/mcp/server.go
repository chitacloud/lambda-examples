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

		if tool := s.FindTool(toolName); tool != nil {
			responseData, err = tool.Handler(r, req.Params.Arguments)
			if err != nil {
				fmt.Printf("Error calling tool %s: %s\n", toolName, err.Error())
			} else {
				// If the response is a slice, we don't need to wrap it.
				// The Response function will handle streaming it.
				val := reflect.ValueOf(responseData)
				if val.Kind() != reflect.Slice {
					// According to JSON-RPC 2.0, we should use 'result' to contain the response content
					unstructuredBytes, err := json.Marshal(responseData)
					if err != nil {
						return nil, err
					}

					responseData = map[string]any{
						"content": []map[string]any{
							{
								"type": "text",
								"text": string(unstructuredBytes),
							},
						},
						"structuredContent": responseData,
					}
				}
			}
		} else {
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

	return Response(mcpInfo, responseData, err)
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

// SetSSEHeaders sets standard Server-Sent Events headers
func SetSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}

// HandleInitialize creates the initialize response data
func (s *Server) HandleInitialize() map[string]interface{} {
	return map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": true,
			},
		},
		"serverInfo": map[string]interface{}{
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
}

func InitHttp(r *http.Request, w http.ResponseWriter, req MCPRequest) (MCPInfo, error) {

	// Set CORS headers to allow MCP Inspector to connect
	SetCORSHeaders(w)

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return MCPInfo{IsPreflight: true}, nil
	}

	// Set SSE headers
	SetSSEHeaders(w)

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

	return MCPInfo{Method: method, RequestID: requestID, IsPreflight: false}, nil

}

func Response(mcpInfo MCPInfo, responseData any, err error) (io.ReadCloser, error) {
	// Use reflection to check if responseData is a slice
	val := reflect.ValueOf(responseData)

	// Format as SSE
	var buffer strings.Builder

	if val.Kind() == reflect.Slice {
		// If it's a slice, iterate and send each element as a separate event
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i).Interface()
			// Per instructions, wrap each element in a single-element slice
			singleElementSlice := []any{elem}

			responseBody, err := FormatMCPServerResponse(mcpInfo.RequestID, mcpInfo.Method, singleElementSlice, err)
			if err != nil {
				// If an error occurs formatting one element, we can decide how to handle it.
				// For now, we'll return the error, stopping the stream.
				return nil, fmt.Errorf("failed to marshal response for slice element %d: %w", i, err)
			}

			// Add data for the current element
			buffer.WriteString("data: ")
			buffer.Write(responseBody)
			buffer.WriteString("\n\n")
		}
	} else {
		// If it's not a slice, handle as a single response
		responseBody, err := FormatMCPServerResponse(mcpInfo.RequestID, mcpInfo.Method, responseData, err)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}

		// Add data
		buffer.WriteString("data: ")
		buffer.Write(responseBody)
		buffer.WriteString("\n\n")
	}

	return io.NopCloser(strings.NewReader(buffer.String())), nil
}
