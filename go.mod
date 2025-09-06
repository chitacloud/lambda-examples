module github.com/chitacloud/lambda-examples

go 1.24.2

replace github.com/chitacloud/lambda-examples/chitacloud-utils => ./chitacloud-utils

replace github.com/chitacloud/lambda-examples/mcp-hour => ./mcp-hour

require (
	github.com/chitacloud/lambda-examples/chitacloud-utils v0.0.0-00010101000000-000000000000
	github.com/chitacloud/lambda-examples/mcp-hour v0.0.0-00010101000000-000000000000
)

require github.com/goccy/go-json v0.10.4 // indirect

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fredyk/westack-go/lambdas v1.0.0-rc01 // indirect
	github.com/getkin/kin-openapi v0.125.0 // indirect
	github.com/go-openapi/jsonpointer v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/invopop/yaml v0.2.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
