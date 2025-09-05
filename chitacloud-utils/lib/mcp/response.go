// Package mcp provides utilities for creating Model Context Protocol (MCP) servers
package mcp

import (
	"encoding/json"
)

const (
	ErrUnkown = -32001
)

/**
https://www.jsonrpc.org/specification

5.1 Error object
When a rpc call encounters an error, the Response Object MUST contain the error member with a value that is a Object with the following members:

code: A Number that indicates the error type that occurred.
This MUST be an integer.

message: A String providing a short description of the error.
The message SHOULD be limited to a concise single sentence.

data: A Primitive or Structured value that contains additional information about the error.
This may be omitted.
The value of this member is defined by the Server (e.g. detailed error information, nested errors etc.).
The error codes from and including -32768 to -32000 are reserved for pre-defined errors. Any code within this range, but not defined explicitly below is reserved for future use. The error codes are nearly the same as those suggested for XML-RPC at the following url: http://xmlrpc-epi.sourceforge.net/specs/rfc.fault_codes.php

| Code | Message | Meaning |
|------|---------|---------|
| -32700 | Parse error | Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text. |
| -32600 | Invalid Request | The JSON sent is not a valid Request object. |
| -32601 | Method not found | The method does not exist / is not available. |
| -32602 | Invalid params | Invalid method parameter(s). |
| -32603 | Internal error | Internal JSON-RPC error. |
| -32000 to -32099 | Server error | Reserved for implementation-defined server-errors. |
The remainder of the space is available for application defined errors.
*/

// JsonRPCError represents a JSON-RPC 2.0 error
type JsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ProgressInfo struct {
	ProgressToken string `json:"progressToken"`
	Progress      int    `json:"progress"`
	Total         int    `json:"total"`
}

// FormatMCPServerResponse formats the response according to JSON-RPC 2.0 / MCP protocol
func FormatMCPServerResponse(id int, method string, streamId string, content any, progressInfo *ProgressInfo, err error) ([]byte, error) {
	responseObj := map[string]any{
		"jsonrpc": "2.0",
	}

	if method == "notifications/progress" && progressInfo != nil {
		responseObj["method"] = method
		// The client expects content to be an array of ToolContent parts.
		responseObj["params"] = map[string]any{
			"streamId":      streamId,
			"content":       content,
			"progressToken": progressInfo.ProgressToken,
			"progress":      progressInfo.Progress,
			"total":         progressInfo.Total,
		}
	} else {
		responseObj["id"] = id
		if err != nil {
			responseObj["error"] = JsonRPCError{Code: ErrUnkown, Message: err.Error(), Data: map[string]any{"content": content}}
		} else {
			responseObj["result"] = content
		}
	}

	return json.Marshal(responseObj)
}
