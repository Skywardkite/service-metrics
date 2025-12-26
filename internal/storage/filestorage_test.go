package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestMemStorage_SaveMetrics(t *testing.T) {
	tests := []struct {
		name     string
		gauges   map[string]float64
		counters map[string]int64
		wantErr  bool
	}{
		{
			name: "success",
			gauges: map[string]float64{
				"Alloc": 1.5,
				"Heap":  2.5,
			},
			counters: map[string]int64{
				"PollCount": 10,
			},
			wantErr: false,
		},
		{
			name:     "empty metrics",
			gauges:   map[string]float64{},
			counters: map[string]int64{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			filePath := filepath.Join(dir, "metrics.txt")

			st := &MemStorage{}

			err := st.SaveMetrics(filePath, tt.gauges, tt.counters)
			if tt.wantErr {
				assert.Assert(t, err != nil)
				return
			}
			assert.NilError(t, err)

			gotGauges, gotCounters, err := LoadMetrics(filePath)
			assert.NilError(t, err)

			assert.DeepEqual(t, tt.gauges, gotGauges)
			assert.DeepEqual(t, tt.counters, gotCounters)
		})
	}
}

func TestMemStorage_SaveMetrics_OpenFileError(t *testing.T) {
	st := &MemStorage{}

	// Путь в несуществующую директорию
	filePath := filepath.Join(t.TempDir(), "no_such_dir", "metrics.txt")

	err := st.SaveMetrics(filePath, map[string]float64{}, map[string]int64{})

	assert.Assert(t, err != nil)
}

func TestLoadMetrics(t *testing.T) {
	tests := []struct {
		name         string
		fileContent  string
		wantGauges   map[string]float64
		wantCounters map[string]int64
		wantErr      bool
	}{
		{
			name:         "file does not exist",
			wantGauges:   map[string]float64{},
			wantCounters: map[string]int64{},
			wantErr:      false,
		},
		{
			name: "success load metrics",
			fileContent: strings.Join([]string{
				`{"id":"Alloc","type":"gauge","value":1.5}`,
				`{"id":"PollCount","type":"counter","delta":10}`,
			}, "\n"),
			wantGauges: map[string]float64{
				"Alloc": 1.5,
			},
			wantCounters: map[string]int64{
				"PollCount": 10,
			},
			wantErr: false,
		},
		{
			name: "skip invalid json lines",
			fileContent: strings.Join([]string{
				`{"id":"Alloc","type":"gauge","value":1.5}`,
				`{this is broken json`,
				`{"id":"PollCount","type":"counter","delta":7}`,
			}, "\n"),
			wantGauges: map[string]float64{
				"Alloc": 1.5,
			},
			wantCounters: map[string]int64{
				"PollCount": 7,
			},
			wantErr: false,
		},
		{
			name: "unknown metric type is ignored",
			fileContent: strings.Join([]string{
				`{"id":"Alloc","type":"gauge","value":1.5}`,
				`{"id":"Unknown","type":"weird","value":10}`,
			}, "\n"),
			wantGauges: map[string]float64{
				"Alloc": 1.5,
			},
			wantCounters: map[string]int64{},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string

			if tt.fileContent != "" {
				tmpFile := filepath.Join(t.TempDir(), "metrics.txt")
				err := os.WriteFile(tmpFile, []byte(tt.fileContent), 0644)
				require.NoError(t, err)
				filePath = tmpFile
			} else {
				// Файл не существует
				filePath = filepath.Join(t.TempDir(), "not_exists.txt")
			}

			gauges, counters, err := LoadMetrics(filePath)

			if tt.wantErr {
				assert.Assert(t, err != nil)
				return
			}

			assert.NilError(t, err)
			assert.DeepEqual(t, tt.wantGauges, gauges)
			assert.DeepEqual(t, tt.wantCounters, counters)
		})
	}
}
