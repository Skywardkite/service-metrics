package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/Skywardkite/service-metrics/internal/config/server_config"
	"github.com/Skywardkite/service-metrics/internal/handler"
	mocks "github.com/Skywardkite/service-metrics/internal/handler/example_mocks"
	model "github.com/Skywardkite/service-metrics/internal/model"
	"go.uber.org/zap"
)

func ExampleHandler_GetMetricJSONHandler() {
	h := handler.NewHandler(&mocks.MockService{}, &server_config.Config{}, &mocks.MockStorage{}, &zap.SugaredLogger{}, &mocks.MockAudit{})

	// Запрашиваем Gauge метрику
	metric := model.Metrics{ID: "requests", MType: model.Gauge}
	data, _ := json.Marshal(metric)
	req := httptest.NewRequest(http.MethodPost, "/get-json", bytes.NewReader(data))
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()

	h.GetMetricJSONHandler(w, req)
	resp := w.Result()

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))

	var respMetric model.Metrics
	json.NewDecoder(resp.Body).Decode(&respMetric)
	fmt.Println("Metric ID:", respMetric.ID)
	fmt.Println("Metric Type:", respMetric.MType)
	fmt.Println("Metric Value:", strconv.FormatFloat(*respMetric.Value, 'f', -1, 64))

	// Output:
	// Status Code: 200
	// Content-Type: application/json
	// Metric ID: requests
	// Metric Type: gauge
	// Metric Value: 42
}
