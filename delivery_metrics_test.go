package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestReportDeliveryUsesStableWriteKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload["idempotency_key"] != "delivery:pkg-1" {
			t.Fatalf("idempotency_key = %v", payload["idempotency_key"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"data":{},"error":null,"metadata":{}}`))
	}))
	defer server.Close()

	client := &MetricsClient{Key: "test", Endpoint: server.URL, HTTP: server.Client(), Sleep: func(time.Duration) {}, MaxRetries: 1}
	if err := client.Report(map[string]any{"type": "counter", "name": "test", "value": 1, "idempotency_key": "delivery:pkg-1"}); err != nil {
		t.Fatal(err)
	}
}
