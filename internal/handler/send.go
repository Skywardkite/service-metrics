package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"

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

    req, err := http.NewRequest(http.MethodPost, url, &buf)
    if err != nil {
        return fmt.Errorf("failed to create request")
    }
    req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("non-OK response status: %d", resp.StatusCode)
    }

    return nil
}