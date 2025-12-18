package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/Skywardkite/service-metrics/internal/agent"
)

func SendBatch(client *retryablehttp.Client, storage *agent.AgentMetrics, serverURL, key string) error {
	metrics := storage.ConvertToBatch()

	// Не отправляем пустые батчи
	if len(metrics) == 0 {
		return nil
	}

	storage.ClearAgentCounter()

	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(jsonData); err != nil {
		return fmt.Errorf("failed to compress data: %w", err)
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("failed to close gzip writer: %w", err)
	}

	url := fmt.Sprintf("%s/updates/", serverURL)

	req, err := retryablehttp.NewRequest("POST", url, &buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	if key != "" {
		hash := SignBody(buf.Bytes(), key)
		req.Header.Set("HashSHA256", hash)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request after retries: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK response status: %d", resp.StatusCode)
	}

	return nil
}
