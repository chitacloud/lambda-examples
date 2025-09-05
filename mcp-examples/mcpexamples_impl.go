package mcpexamples

import (
	"net/http"

	"github.com/chitacloud/lambda-examples/chitacloud-utils/lib/mcp"
	"github.com/getkin/kin-openapi/openapi3"
)

func exampleSliceHandler(r *http.Request, params map[string]any) (any, error) {
	// This handler returns a slice of data to test the streaming functionality.
	data := []map[string]any{
		{"id": 1, "name": "first item"},
		{"id": 2, "name": "second item"},
		{"id": 3, "name": "third item"},
	}
	return data, nil
}

func registerExampleSliceTool(server *mcp.Server) {
	server.RegisterTool(mcp.ToolDescription{
		Name:        "example_slice",
		Description: "An example tool that returns a slice of items to demonstrate streaming.",
		InputSchema: &openapi3.Schema{
			Type: &openapi3.Types{openapi3.TypeObject},
		},
		OutputSchema: &openapi3.Schema{
			Type: &openapi3.Types{openapi3.TypeObject},
			// First item will be the {count: x} of items
			// Followed by the items one by one as single entries
			// oneOf
			OneOf: []*openapi3.SchemaRef{
				{Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeObject}, Properties: map[string]*openapi3.SchemaRef{
					"count": {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeInteger}}}}}},
				{Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeObject}, Properties: map[string]*openapi3.SchemaRef{
					"id":   {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeInteger}}},
					"name": {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}},
				}}},
			},
		},
		Handler: exampleSliceHandler,
	})
}
