package mcpexamples

import (
	"io"
	"net/http"

	chitamcputils "github.com/chitacloud/lambda-examples/chitacloud-utils/chita-mcp-utils"
	"github.com/chitacloud/lambda-examples/chitacloud-utils/lib/mcp"
)

var server *mcp.Server

func init() {
	server = chitamcputils.DefaultServer(chitamcputils.ServerOptions{
		Name:        "MCP Examples",
		Version:     "1.0.0",
		Description: "MCP Examples",
		Debug:       true,
	})

	registerExampleSliceTool(server)
}

func ExamplesHandler(r *http.Request, w http.ResponseWriter, req mcp.MCPRequest) (io.ReadCloser, error) {
	return server.Handle(r, w, req)
}
