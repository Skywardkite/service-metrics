package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/Skywardkite/service-metrics/internal/config/server_config"
	"github.com/Skywardkite/service-metrics/internal/handler"
	mocks "github.com/Skywardkite/service-metrics/internal/handler/example_mocks"
	model "github.com/Skywardkite/service-metrics/internal/model"
	pointers "github.com/Skywardkite/service-metrics/internal/utils"
	"go.uber.org/zap"
)

func ExampleHandler_UpdateMetricsBatchJSONHandler() {
	h := handler.NewHandler(&mocks.MockService{}, &server_config.Config{}, &mocks.MockStorage{}, &zap.SugaredLogger{}, &mocks.MockAudit{})

	// Создаём батч из двух метрик: Gauge и Counter
	gv := pointers.To(float64(42.0))
	cd := int64(7)
	batch := []model.Metrics{
		{ID: "requests", MType: model.Gauge, Value: gv},
		{ID: "tasks_processed", MType: model.Counter, Delta: &cd},
	}

	data, _ := json.Marshal(batch)
	req := httptest.NewRequest(http.MethodPost, "/update-batch-json", bytes.NewReader(data))
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()

	h.UpdateMetricsBatchJSONHandler(w, req)
	resp := w.Result()

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	fmt.Println("Response body:", respBody)

	// Output:
	// Status Code: 200
	// Content-Type: application/json
	// Response body: map[]
}
