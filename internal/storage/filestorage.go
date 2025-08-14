package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	model "github.com/Skywardkite/service-metrics/internal/model"
)

func SaveMetrics(filePath string, gauges map[string]float64, counters map[string]int64) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
    if err != nil {
        return err
    }
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Записываем все gauge метрики
	for id, value := range gauges {
		data, _ := json.Marshal(struct {
			ID    string  `json:"id"`
			Type  string  `json:"type"`
			Value float64 `json:"value"`
		}{
			ID:    id,
			Type:  "gauge",
			Value: value,
		})
		if _, err := writer.Write(append(data, '\n')); err != nil {
			os.Remove(filePath)
			return err
		}
	}

	// Записываем все counter метрики
	for id, delta := range counters {
		data, _ := json.Marshal(struct {
			ID    string `json:"id"`
			Type  string `json:"type"`
			Delta int64  `json:"delta"`
		}{
			ID:    id,
			Type:  "counter",
			Delta: delta,
		})
		if _, err := writer.Write(append(data, '\n')); err != nil {
			os.Remove(filePath)
			return err
		}
	}

	if err := writer.Flush(); err != nil {
        os.Remove(filePath)
        return fmt.Errorf("flush error: %w", err)
    }

    // 5. Атомарная замена файла
    if err := os.Rename(filePath, filePath); err != nil {
        os.Remove(filePath)
        return fmt.Errorf("rename error: %w", err)
    }

	return nil
}

func LoadMetrics(filePath string) (map[string]float64, map[string]int64, error) {
	gauges := make(map[string]float64)
	counters := make(map[string]int64)

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return gauges, counters, nil
		}
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var m model.Metrics
		if err := json.Unmarshal(scanner.Bytes(), &m); err != nil {
			continue // Пропускаем некорректные строки
		}

		switch m.MType {
		case "gauge":
			gauges[m.ID] = *m.Value
		case "counter":
			counters[m.ID] = *m.Delta 
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return gauges, counters, nil
}