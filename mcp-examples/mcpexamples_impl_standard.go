package mcpexamples

import (
	"github.com/chitacloud/lambda-examples/chitacloud-utils/lib/mcp"
	"github.com/getkin/kin-openapi/openapi3"
)

func registerStandardSliceTool(server *mcp.Server) {
	server.RegisterTool(mcp.ToolDescription{
		Name:        "standard_slice_tool",
		Description: "An example tool that returns a slice of items to demonstrate standard streaming.",
		InputSchema: &openapi3.Schema{
			Type: &openapi3.Types{openapi3.TypeObject},
		},
		OutputSchema: &openapi3.Schema{
			Type: &openapi3.Types{openapi3.TypeObject},
			Properties: map[string]*openapi3.SchemaRef{
				"items": {
					Value: &openapi3.Schema{
						Type: &openapi3.Types{openapi3.TypeArray},
						Items: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: &openapi3.Types{openapi3.TypeObject},
								Properties: map[string]*openapi3.SchemaRef{
									"id":   {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeInteger}}},
									"name": {Value: &openapi3.Schema{Type: &openapi3.Types{openapi3.TypeString}}},
								},
								Required: []string{"id", "name"},
							},
						},
					},
				},
			},
		},
		Handler: exampleSliceHandler,
	})
}
