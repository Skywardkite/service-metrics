package agent

import (
	model "github.com/Skywardkite/service-metrics/internal/model"
)

func (s *AgentMetrics) ConvertToBatch() []model.Metrics {
	var batch []model.Metrics

	for name, value := range s.Gauge {
		batch = append(batch, model.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		})
	}

	for name, value := range s.Counter {
		batch = append(batch, model.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &value,
		})
	}

	return batch
}