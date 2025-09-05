package mcpexamples

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chitacloud/lambda-examples/chitacloud-utils/lib/mcp"
	"github.com/fredyk/westack-go/v2/lambdas"
	"github.com/stretchr/testify/assert"
)

// MCPTestResponse is a helper struct for parsing MCP streaming responses in tests.
type MCPTestResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		StreamID string `json:"streamId"`
		Content  any    `json:"content"`
	} `json:"params"`
	Result any `json:"result"`
	ID     int `json:"id"`
}

func executeMCPToolCall(t *testing.T, server *mcp.Server, toolName string, arguments map[string]any) *http.Response {
	t.Helper()

	req := mcp.MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params: mcp.MCPRequestParams{
			Name:      toolName,
			Arguments: arguments,
		},
	}

	body, err := json.Marshal(req)
	assert.NoError(t, err, "Failed to marshal request")

	httpReq := httptest.NewRequest("POST", "/", strings.NewReader(string(body)))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	req.LambdaRequest = lambdas.LambdaRequest{Payload: body}

	respBody, err := server.Handle(httpReq, w, req)
	assert.NoError(t, err, "Handler returned an error")
	defer respBody.Close()

	_, err = io.Copy(w.Body, respBody)
	assert.NoError(t, err, "Failed to copy response body")

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")

	return resp
}

func TestRawSliceStreaming(t *testing.T) {
	// Create a new server
	server := mcp.NewServer("test-server", "1.0", "Test Server")
	registerExampleSliceTool(server)

	// Execute the tool call using the helper
	resp := executeMCPToolCall(t, server, "example_slice", map[string]any{})

	// Verify the SSE stream
	scanner := bufio.NewScanner(resp.Body)
	defer resp.Body.Close()

	var eventCount int
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data:") {
			eventCount++
		}
	}

	assert.NoError(t, scanner.Err(), "Error reading response body")

	// The example_slice tool returns a slice with 3 items
	assert.Equal(t, 3, eventCount, "Expected 3 SSE events for raw streaming")
}
