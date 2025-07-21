package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Skywardkite/service-metrics/internal/agent"
)

func SendMetrics(storage *agent.AgentMetrics, port string) {
    client := &http.Client{}
    gauges, counters := storage.GetAgentMetrics()

    if !strings.HasPrefix(port, "http://") && !strings.HasPrefix(port, "https://") {
        port = "http://" + port
    }
    
    for name, value := range gauges {
        url := fmt.Sprintf("%s/update/gauge/%s/%f", port, name, value)
        sendPlainPost(client, url)
    }

    for name, delta := range counters {
        url := fmt.Sprintf("%s/update/counter/%s/%d", port, name, delta)
        sendPlainPost(client, url)
    }
}

func sendPlainPost(client *http.Client, url string) error {
    req, err := http.NewRequest(http.MethodPost, url, nil)
    if err != nil {
        return errors.New("failed to create request")
    }
    req.Header.Set("Content-Type", "text/plain")

    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return errors.New("non-OK response status")
    }

    return nil
}