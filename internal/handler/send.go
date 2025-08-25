package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/Skywardkite/service-metrics/internal/agent"
	model "github.com/Skywardkite/service-metrics/internal/model"
)

func SendMetrics(client *http.Client, storage *agent.AgentMetrics, url string) {
    gauges, counters := storage.GetAgentMetrics()

    for name, value := range gauges {
        sendPlainPost(client, url, model.Metrics{
            ID: name,
            MType: model.Gauge,
            Value: &value,
        })
    }
        
    for name, delta := range counters {
        sendPlainPost(client, url, model.Metrics{
            ID: name,
            MType: model.Counter,
            Delta: &delta,
        })
    }

    //После отправки метрик обнуляем счетчик сбора
    storage.ClearAgentCounter()
}

func sendPlainPost(client *http.Client, url string, metric model.Metrics) error {
    jsonData, err := json.Marshal(metric)
    if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
    }

    // Сжимаем данные
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(jsonData); err != nil {
		return fmt.Errorf("failed to compress data: %w", err)
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("failed to close gzip writer: %w", err)
	}

    //повторы при проблемах с отправкой метрик
    delays := []time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

    var lastErr error
    for attempt := 0; attempt <= len(delays); attempt++ {
        bodyCopy := bytes.NewBuffer(buf.Bytes())

        req, err := http.NewRequest(http.MethodPost, url, bodyCopy)
        if err != nil {
            return fmt.Errorf("failed to create request: %w", err)
        }
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Content-Encoding", "gzip")
        req.Header.Set("Accept-Encoding", "gzip")

        resp, err := client.Do(req)
        if err == nil && resp.StatusCode == http.StatusOK {
            resp.Body.Close()
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

func isRetriableError(err error) bool {
    if ne, ok := err.(net.Error); ok && (ne.Temporary() || ne.Timeout()) {
        return true
    }
    if errors.Is(err, io.ErrUnexpectedEOF) {
        return true
    }
    
    return false
}