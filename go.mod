module github.com/chitacloud/lambda-examples

go 1.24.2

replace github.com/chitacloud/lambda-examples/chitacloud-utils => ./chitacloud-utils

replace github.com/chitacloud/lambda-examples/mcp-hour => ./mcp-hour

require (
	github.com/chitacloud/lambda-examples/chitacloud-utils v0.0.0-00010101000000-000000000000
	github.com/chitacloud/lambda-examples/mcp-hour v0.0.0-00010101000000-000000000000
)

require (
	github.com/fredyk/westack-go/v2/lambdas v0.0.0-20250904092258-4f6557da7f99 // indirect
	github.com/goccy/go-json v0.10.4 // indirect
)
