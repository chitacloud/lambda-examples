module github.com/chitacloud/example-lambdas/mcp-hour

go 1.24.2

replace github.com/chitacloud/example-lambdas/chita-utils => ../chitacloud-utils

require github.com/chitacloud/example-lambdas/chita-utils v0.0.0-00010101000000-000000000000

require (
	github.com/fredyk/westack-go/v2/lambdas v0.0.0-20250904092258-4f6557da7f99 // indirect
	github.com/goccy/go-json v0.10.4 // indirect
)
