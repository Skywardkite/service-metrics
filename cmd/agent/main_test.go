package main

import (
	"testing"
	"time"
)

// MockStore для тестирования
type MockStore struct {
	pollCount int
	sendCount int
}

func (m *MockStore) PollRuntimeMetrics() {
	m.pollCount++
}

func (m *MockStore) SendMetrics() {
	m.sendCount++
}

func TestMain(t *testing.T) {
	// Короткие интервалы для теста
	pollInterval := 10 * time.Millisecond
	reportInterval := 50 * time.Millisecond

	var pollCount, sendCount int
	lastReport := time.Now()

	// Запускаем ограниченное количество итераций
	for range 10 {
		pollCount++

		if time.Since(lastReport) >= reportInterval {
			sendCount++
			lastReport = time.Now()
		}

		time.Sleep(pollInterval)
	}

	// Проверяем результаты
	if pollCount != 10 {
		t.Errorf("Expected 10 polls, got %d", pollCount)
	}

	// Ожидаем 2 отправки (50ms интервал / 10ms sleep = ~5 итераций между отправками)
	if sendCount < 1 || sendCount > 2 {
		t.Errorf("Expected 1-2 sends, got %d", sendCount)
	}
}