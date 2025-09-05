package mcpexamples

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chitacloud/lambda-examples/chitacloud-utils/lib/mcp"
	"github.com/fredyk/westack-go/v2/lambdas"
)

func TestExampleSliceStreaming(t *testing.T) {
	// Create a new server
	server := mcp.NewServer("test-server", "1.0", "Test Server")
	registerExampleSliceTool(server)

	// Create a request for the example_slice tool
	req := mcp.MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params: mcp.MCPRequestParams{
			Name:      "example_slice",
			Arguments: map[string]any{},
		},
	}

	// Marshal the request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Create a mock HTTP request
	httpReq := httptest.NewRequest("POST", "/", strings.NewReader(string(body)))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// The server's Handle method expects an MCPRequest, which includes the raw payload.
	req.LambdaRequest = lambdas.LambdaRequest{Payload: body}

	respBody, err := server.Handle(httpReq, w, req)
	if err != nil {
		t.Fatalf("Handler returned an error: %v", err)
	}
	defer respBody.Close()
	io.Copy(w.Body, respBody)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Verify the SSE stream
	scanner := bufio.NewScanner(resp.Body)
	defer resp.Body.Close()

	var eventCount int
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data:") {
			eventCount++
			// You can add more detailed checks here, e.g., unmarshal the data
			// and verify its content.
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	// The example_slice tool returns a slice with 3 items, plus the count
	expectedEventCount := 4
	if eventCount != expectedEventCount {
		t.Errorf("Expected %d SSE events, got %d", expectedEventCount, eventCount)
	}
}