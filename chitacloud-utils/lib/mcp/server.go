// Package mcp provides utilities for creating Model Context Protocol (MCP) servers
package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Server represents an MCP protocol server
type Server struct {
	Name           string
	Version        string
	Description    string
	Tools          []ToolDescription
	DefaultHandler func(params map[string]any) (map[string]any, error)
}

// NewServer creates a new MCP server with the given parameters
func NewServer(name, version, description string) *Server {
	return &Server{
		Name:        name,
		Version:     version,
		Description: description,
		Tools:       []ToolDescription{},
	}
}

// RegisterTool adds a tool to the server's available tools
func (s *Server) RegisterTool(tool ToolDescription) {
	s.Tools = append(s.Tools, tool)
}

func (s *Server) SetDefaultHandler(handler func(params map[string]any) (map[string]any, error)) {
	s.DefaultHandler = handler
}

func (s *Server) Handle(w http.ResponseWriter, r *http.Request, req MCPRequest, mcpInfo MCPInfo) (io.ReadCloser, error) {
	// Prepare the response based on path
	var responseData map[string]any
	var err error

	// Handle different MCP protocol paths
	switch mcpInfo.Method {
	case "initialize":
		// Initialize request - return server capabilities
		responseData = s.HandleInitialize()
		fmt.Println("Sending initialize response")

	case "tools/list":
		// List tools request
		responseData = s.HandleTools()
		fmt.Println("Sending tools list response")

	case "tools/call":
		toolName := req.Params["name"]

		if toolName, ok := toolName.(string); ok {
			tool := s.FindTool(toolName)
			if tool != nil {
				responseData, err = tool.Handler(req.Params)
				if err != nil {
					responseData = map[string]any{"error": err.Error()}
				}
			} else {
				responseData = map[string]any{"error": "tool not found"}
			}
		} else {
			responseData = map[string]any{"error": "tool name must be a string"}
		}
	default:
		// Default path - for compatibility with legacy clients
		responseData, err = s.DefaultHandler(req.Params)
		fmt.Println("Sending default path response")
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

	// // Extract method from path
	// pathMethod := GetMethodFromPath(r.URL.Path)
	// fmt.Printf("Path method: %s\n", pathMethod)
	paramsB, err := json.Marshal(req.Params)
	if err != nil {
		return MCPInfo{}, err
	}
	fmt.Printf("Params: %s\n", string(paramsB))

	if id, ok := req.Params["requestId"]; ok {
		if idStr, ok := id.(string); ok {
			requestID, _ = strconv.Atoi(idStr)
		} else if idInt, ok := id.(int); ok {
			requestID = idInt
		}
	}

	// Use request method if available, otherwise use path-derived method
	method := req.Method
	if method == "" {
		// method = pathMethod
		// if method == "" {
		method = "response"
		// }
	}

	return MCPInfo{Method: method, RequestID: requestID, IsPreflight: false}, nil

}

func Response(mcpInfo MCPInfo, responseData interface{}, err error) (io.ReadCloser, error) {

	// Format as SSE
	var buffer strings.Builder

	if err != nil {
		responseData = map[string]any{"error": err.Error()}
	} else {

		// Format as JSON-RPC 2.0 response for MCP
		responseBody, err := FormatMCPServerResponse(mcpInfo.RequestID, mcpInfo.Method, responseData)

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
