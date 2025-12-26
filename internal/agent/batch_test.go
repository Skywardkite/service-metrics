package agent

import (
	"testing"

	"gotest.tools/assert"
)

func TestAgentMetrics_ConvertToBatch(t *testing.T) {
	am := &AgentMetrics{
		Gauge: map[string]float64{
			"Alloc": 123.45,
			"Heap":  678.9,
		},
		Counter: map[string]int64{
			"PollCount": 10,
			"Errors":    2,
		},
	}

	batch := am.ConvertToBatch()

	assert.Equal(t, len(batch), len(am.Gauge)+len(am.Counter))

	gauges := map[string]float64{}
	counters := map[string]int64{}
	for _, m := range batch {
		switch m.MType {
		case "gauge":
			gauges[m.ID] = *m.Value
		case "counter":
			counters[m.ID] = *m.Delta
		default:
			t.Fatalf("unexpected metric type: %s", m.MType)
		}
	}

	assert.DeepEqual(t, am.Gauge, gauges)
	assert.DeepEqual(t, am.Counter, counters)
}
