package mcp_hour

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/chitacloud/example-lambdas/chita-utils/lib/mcp"
)

// testResponseWriter is a mock http.ResponseWriter for testing
type testResponseWriter struct {
	header     http.Header
	statuscode int
	body       strings.Builder
}

func (w *testResponseWriter) Header() http.Header {
	return w.header
}

func (w *testResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *testResponseWriter) WriteHeader(statusCode int) {
	w.statuscode = statusCode
}

// parseSSEEvent parses an SSE event data line into a HourResponse
func parseSSEEvent(dataLine string) (HourResponse, error) {
	// SSE data lines are prefixed with "data: "
	dataLine = strings.TrimPrefix(dataLine, "data: ")

	// Try parsing JSON-RPC 2.0 response format where hour data is in 'result.structuredContent' field
	var jsonRpcResp struct {
		JsonRPC string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  struct {
			StructuredContent string `json:"structuredContent"`
		} `json:"result"`
	}

	err := json.Unmarshal([]byte(dataLine), &jsonRpcResp)
	if err == nil && jsonRpcResp.Result.StructuredContent != "" {
		var hourResp HourResponse
		err := json.Unmarshal([]byte(jsonRpcResp.Result.StructuredContent), &hourResp)
		if err != nil {
			return HourResponse{}, fmt.Errorf("failed to unmarshal structuredContent: %w", err)
		}
		return hourResp, nil
	}

	// Try direct format (for backward compatibility)
	var directResp HourResponse
	err = json.Unmarshal([]byte(dataLine), &directResp)
	if err == nil && (directResp.Hour > 0 || directResp.AmPm != "") {
		return directResp, nil
	}

	return HourResponse{}, fmt.Errorf("failed to parse SSE event: %s", dataLine)
}

func TestHandler(t *testing.T) {
	// Create a dummy request
	req, err := http.NewRequest("GET", "/sse", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a fake ResponseWriter for testing
	w := &testResponseWriter{header: make(http.Header)}

	// Call the handler which returns an io.ReadCloser
	resp, err := Handler(req, w, mcp.MCPRequest{})
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	// Verify we got an io.ReadCloser back
	if resp == nil {
		t.Fatal("Handler returned nil response")
	}

	// Make sure we close the reader when done
	defer resp.Close()

	// Parse the response as a HourResponse using our parseSSEEvent function
	scanner := bufio.NewScanner(resp)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 {
		t.Fatal("No data lines found in SSE response")
	}

	// Parse the first data line as a HourResponse
	response, err := parseSSEEvent(lines[0])
	if err != nil {
		t.Fatalf("Failed to parse SSE event: %v\nLine content: %s", err, lines[0])
	}

	// Verify the response structure
	if response.Hour < 1 || response.Hour > 12 {
		t.Errorf("Hour should be between 1 and 12, got %d", response.Hour)
	}

	if response.AmPm != "AM" && response.AmPm != "PM" {
		t.Errorf("AmPm should be either AM or PM, got %s", response.AmPm)
	}

	// Check if the message contains the hour and AM/PM
	expectedMessagePart := "Current hour is "
	if response.Message[:len(expectedMessagePart)] != expectedMessagePart {
		t.Errorf("Message should start with '%s', got '%s'", expectedMessagePart, response.Message)
	}

	// Ensure the currentTime string is a valid RFC3339 timestamp
	if _, err := time.Parse(time.RFC3339, response.CurrentTime); err != nil {
		t.Errorf("Expected currentTime to be a valid RFC3339 timestamp, but got error: %v", err)
	}
}
