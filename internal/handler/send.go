package handler

import (
	"bytes"
	"compress/gzip"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/Skywardkite/service-metrics/internal/agent"
	"github.com/Skywardkite/service-metrics/internal/crypto"
	model "github.com/Skywardkite/service-metrics/internal/model"
)

// SendMetrics отправляет метрики на сервер.
func SendMetrics(client *retryablehttp.Client, storage *agent.AgentMetrics, url, key string, pubKey *rsa.PublicKey) {
	gauges, counters := storage.GetAgentMetrics()

	for name, value := range gauges {
		sendPlainPost(client, url, key, model.Metrics{
			ID:    name,
			MType: model.Gauge,
			Value: &value,
		}, pubKey)
	}

	for name, delta := range counters {
		sendPlainPost(client, url, key, model.Metrics{
			ID:    name,
			MType: model.Counter,
			Delta: &delta,
		}, pubKey)
	}

	// После отправки метрик обнуляем счетчик сбора.
	storage.ClearAgentCounter()
}

func sendPlainPost(client *retryablehttp.Client, url, key string, metric model.Metrics, pubKey *rsa.PublicKey) error {
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

	bodyBytes := buf.Bytes()

	// Шифруем данные
	if pubKey != nil {
		encrypted, err := crypto.Encrypt(pubKey, bodyBytes)
		if err != nil {
			return fmt.Errorf("failed to encrypt body: %w", err)
		}
		bodyBytes = encrypted
	}

	req, err := retryablehttp.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
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
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK response status: %d", resp.StatusCode)
	}

	return nil
}
