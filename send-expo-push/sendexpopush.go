package sendexpopush

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Request represents the input for sending an Expo push notification.
type Request struct {
	// Expo push token, e.g. "ExpoPushToken[xxxxxxxxxxxxxxxxxxxxxx]"
	To string `json:"to"`
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
	Data  map[string]any `json:"data,omitempty"`
	Sound string `json:"sound,omitempty"`
	TTL   *int   `json:"ttl,omitempty"`
	// normal | high
	Priority string `json:"priority,omitempty"`
}

// Response is a simplified response from Expo Push service.
type Response struct {
	Status  string            `json:"status"`
	ID      string            `json:"id,omitempty"`
	Message string            `json:"message,omitempty"`
	Details map[string]any    `json:"details,omitempty"`
	// Raw captures the raw response for debugging purposes
	Raw     map[string]any    `json:"-"`
}

var (
	defaultHTTPClient = &http.Client{Timeout: 10 * time.Second}
)

// getBaseURL returns the Expo Push API base URL. Overridable via EXPO_PUSH_BASE_URL for tests.
func getBaseURL() string {
	if v := os.Getenv("EXPO_PUSH_BASE_URL"); v != "" {
		return strings.TrimRight(v, "/")
	}
	return "https://exp.host"
}

// validateToken performs a basic sanity check for Expo push tokens.
func validateToken(token string) error {
	if token == "" {
		return errors.New("to (Expo push token) is required")
	}
	if !(strings.HasPrefix(token, "ExpoPushToken[") || strings.HasPrefix(token, "ExponentPushToken[")) {
		// Don't hard fail for non-standard prefixes, but warn via error for better DX
		return fmt.Errorf("unexpected push token format: %s", token)
	}
	return nil
}

// SendExpoPush sends a push notification using Expo's REST API.
//
// ChitaCloud handler style: Accepts a typed Request and returns a typed Response and error.
func SendExpoPush(req Request) (Response, error) {
	// Validate token lightly; if format mismatch we still allow sending after warning
	if err := validateToken(req.To); err != nil {
		// If token is empty or clearly invalid, abort. For prefix mismatch, we'll still attempt.
		if req.To == "" {
			return Response{Status: "error", Message: err.Error()}, err
		}
	}

	payload := map[string]any{
		"to":    req.To,
		"title": req.Title,
		"body":  req.Body,
	}
	if len(req.Data) > 0 {
		payload["data"] = req.Data
	}
	if req.Sound != "" {
		payload["sound"] = req.Sound
	}
	if req.TTL != nil {
		payload["ttl"] = *req.TTL
	}
	if req.Priority != "" {
		payload["priority"] = req.Priority
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return Response{Status: "error", Message: "failed to encode payload"}, err
	}

	url := getBaseURL() + "/--/api/v2/push/send"
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return Response{Status: "error", Message: "failed to create http request"}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	if token := os.Getenv("EXPO_ACCESS_TOKEN"); token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := defaultHTTPClient.Do(httpReq)
	if err != nil {
		return Response{Status: "error", Message: "failed to call expo push API"}, err
	}
	defer resp.Body.Close()

	var out map[string]any
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&out); err != nil {
		return Response{Status: "error", Message: fmt.Sprintf("failed to decode response: %v", err)}, err
	}

	// Expo returns 200 with a data object or 200 with error status.
	// Structure typically: { "data": { "status": "ok", "id": "..." } }
	// or { "data": { "status": "error", "message": "...", "details": {...} } }
	data, _ := out["data"].(map[string]any)
	if data == nil {
		// Sometimes it may be an array; handle first element
		if arr, ok := out["data"].([]any); ok && len(arr) > 0 {
			data, _ = arr[0].(map[string]any)
		}
	}

	if data == nil {
		// Fallback: return raw
		return Response{Status: "error", Message: "unexpected response shape", Details: out, Raw: out}, fmt.Errorf("unexpected response shape")
	}

	status, _ := data["status"].(string)
	id, _ := data["id"].(string)
	message, _ := data["message"].(string)
	var details map[string]any
	if d, ok := data["details"].(map[string]any); ok {
		details = d
	}

	res := Response{Status: status, ID: id, Message: message, Details: details, Raw: out}
	if status != "ok" {
		if message == "" {
			message = fmt.Sprintf("expo push returned status %s", status)
		}
		return res, fmt.Errorf("%s", message)
	}
	return res, nil
}
