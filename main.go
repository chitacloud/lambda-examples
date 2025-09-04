package main

import (
	"net/http"

	"github.com/chitacloud/lambda-examples/chitacloud-utils/lib/mcp"
	mcp_hour "github.com/chitacloud/lambda-examples/mcp-hour"
)

func main() {
	mcp_hour.Handler(&http.Request{}, nil, mcp.MCPRequest{})
}
