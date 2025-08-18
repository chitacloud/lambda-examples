package sendexpopush

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSendExpoPush_Success(t *testing.T) {
	// Mock Expo API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/--/api/v2/push/send" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"status": "ok",
				"id":     "abcdef-123456",
			},
		})
	}))
	defer ts.Close()

	os.Setenv("EXPO_PUSH_BASE_URL", ts.URL)
	defer os.Unsetenv("EXPO_PUSH_BASE_URL")

	res, err := SendExpoPush(Request{
		To:    "ExpoPushToken[xxxxxxxxxxxxxxxxxxxxxx]",
		Title: "Hello",
		Body:  "World",
		Data:  map[string]any{"k": "v"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != "ok" {
		t.Fatalf("expected status ok, got %s", res.Status)
	}
	if res.ID == "" {
		t.Fatalf("expected an id in response")
	}
}

func TestSendExpoPush_ErrorFromAPI(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"status":  "error",
				"message": "DeviceNotRegistered",
				"details": map[string]any{"error": "DeviceNotRegistered"},
			},
		})
	}))
	defer ts.Close()

	os.Setenv("EXPO_PUSH_BASE_URL", ts.URL)
	defer os.Unsetenv("EXPO_PUSH_BASE_URL")

	_, err := SendExpoPush(Request{To: "ExpoPushToken[bad]"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestSendExpoPush_EmptyToken(t *testing.T) {
	_, err := SendExpoPush(Request{})
	if err == nil {
		t.Fatalf("expected error for empty token")
	}
}

func TestSendExpoPush_Integration(t *testing.T) {
	// Only run if a real token is provided
	pushToken := os.Getenv("EXPO_INTEGRATION_PUSH_TOKEN")
	if pushToken == "" {
		t.Skip("EXPO_INTEGRATION_PUSH_TOKEN not set; skipping integration test")
	}
	// Optional: EXPO_ACCESS_TOKEN may be needed for some accounts
	res, err := SendExpoPush(Request{
		To:    pushToken,
		Title: "Test from CI",
		Body:  "Hello from integration test",
	})
	if err != nil {
		t.Fatalf("integration send failed: %v (res=%+v)", err, res)
	}
	if res.Status != "ok" {
		t.Fatalf("expected ok, got %s", res.Status)
	}
}
