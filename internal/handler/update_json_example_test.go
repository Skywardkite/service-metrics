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
	pointers "github.com/Skywardkite/service-metrics/internal/utils"
	"go.uber.org/zap"
)

func ExampleHandler_UpdateJSONHandler() {
	h := handler.NewHandler(&mocks.MockService{}, &server_config.Config{}, &mocks.MockStorage{}, &zap.SugaredLogger{}, &mocks.MockAudit{})

	m := model.Metrics{
		ID:    "requests",
		MType: model.Gauge,
		Value: pointers.To(float64(42.0)),
	}

	data, _ := json.Marshal(m)
	req := httptest.NewRequest(http.MethodPost, "/update-json", bytes.NewReader(data))
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()

	h.UpdateJSONHandler(w, req)
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
