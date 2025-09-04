package chitamcputils

import (
	"sync"

	"github.com/chitacloud/lambda-examples/chitacloud-utils/lib/mcp"
)

type ServerOptions struct {
	Name        string
	Version     string
	Description string
}

// CreateMCPServer initializes and configures an MCP server for our hour service
func CreateMCPServer(options ServerOptions) *mcp.Server {
	server := mcp.NewServer(options.Name, options.Version, options.Description)
	return server
}

// sync.Once
var once sync.Once
var server *mcp.Server

func DefaultServer(options ServerOptions) *mcp.Server {
	once.Do(func() {
		server = CreateMCPServer(options)
	})
	return server
}
