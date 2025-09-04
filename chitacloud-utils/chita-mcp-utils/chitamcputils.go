package chitamcputils

import "github.com/chitacloud/chita-utils/lib/mcp"

// CreateMCPServer initializes and configures an MCP server for our hour service
func CreateMCPServer() *mcp.Server {
	server := mcp.NewServer("HourMCP", "1.0.0", "MCP server that provides current hour information")
	return server
}
