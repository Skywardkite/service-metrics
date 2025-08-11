package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Skywardkite/service-metrics/internal/agent"
	model "github.com/Skywardkite/service-metrics/internal/model"
)

func SendMetrics(client *http.Client, storage *agent.AgentMetrics, url string) error {
    gauges, counters := storage.GetAgentMetrics()
	var metrics []model.Metrics
	for name, value := range gauges {
		v := value
		metrics = append(metrics, model.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &v,
		})
	}
	for name, delta := range counters {
		d := delta
		metrics = append(metrics, model.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &d,
		})
	}

    for _, metric := range metrics {
        jsonData, err := json.Marshal(metric)
        if err != nil {
		    return fmt.Errorf("failed to marshal metrics: %w", err)
        }

        req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
        if err != nil {
		    return fmt.Errorf("failed to create request")
        }
        req.Header.Set("Content-Type", "application/json")

        resp, err := client.Do(req)
        if err != nil {
		    return fmt.Errorf("request failed: %w", err)
        }
        resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
		    return fmt.Errorf("non-OK response status: %d", resp.StatusCode)
        }
    }

	storage.ClearAgentCounter()

	return nil
}