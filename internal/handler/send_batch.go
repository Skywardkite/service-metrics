package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Skywardkite/service-metrics/internal/agent"
)

func SendBatch(client *http.Client, storage *agent.AgentMetrics, serverURL string) error {
	metrics := storage.ConvertToBatch()
	
	// Не отправляем пустые батчи
	if len(metrics) == 0 {
		return nil
	}

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

	//повторы при проблемах с отправкой метрик
    delays := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

    var lastErr error
    for attempt := 0; attempt <= len(delays); attempt++ {
        bodyCopy := bytes.NewBuffer(buf.Bytes())

		req, err := http.NewRequest("POST", url, bodyCopy)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")

		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
            resp.Body.Close()

			//После отправки метрик обнуляем счетчик сбора
    		storage.ClearAgentCounter()
            return nil
        }

		if err != nil {
            if !isRetriableError(err) {
                return fmt.Errorf("non-retriable request error: %w", err)
            }
            lastErr = err
        } else {
            resp.Body.Close()
            if resp.StatusCode >= 500 {
                lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
            } else {
                return fmt.Errorf("non-OK response: %d", resp.StatusCode)
            }
        }

        if attempt < len(delays) {
            time.Sleep(delays[attempt])
        }
	}

	return fmt.Errorf("all retries failed, last error: %w", lastErr)
}