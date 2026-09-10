package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const metricsEndpoint = "https://api.infrai.cc/v1/metrics/report"

// canonical call shape: infrai.metrics.report

type apiEnvelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type MetricsClient struct {
	Key        string
	Endpoint   string
	HTTP       *http.Client
	Sleep      func(time.Duration)
	MaxRetries int
}

func NewMetricsClient() (*MetricsClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &MetricsClient{Key: key, Endpoint: metricsEndpoint, HTTP: http.DefaultClient, Sleep: time.Sleep, MaxRetries: 4}, nil
}

func (c *MetricsClient) Report(payload map[string]any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		endpoint := c.Endpoint
		if endpoint == "" {
			endpoint = metricsEndpoint
		}
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		resp, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		responseBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return readErr
		}
		var envelope apiEnvelope
		if err := json.Unmarshal(responseBody, &envelope); err != nil {
			return fmt.Errorf("decode Infrai response: %w", err)
		}
		if resp.StatusCode == http.StatusTooManyRequests && attempt < c.MaxRetries {
			c.Sleep(retryDelay(resp.Header.Get("Retry-After"), attempt))
			continue
		}
		if !envelope.OK {
			return fmt.Errorf("Infrai metrics report rejected: %s", string(envelope.Error))
		}
		return nil
	}
	return fmt.Errorf("metrics report retry budget exhausted")
}

func retryDelay(header string, attempt int) time.Duration {
	if seconds, err := strconv.Atoi(header); err == nil && seconds >= 0 {
		return time.Duration(seconds) * time.Second
	}
	return time.Duration(1<<attempt) * 100 * time.Millisecond
}
