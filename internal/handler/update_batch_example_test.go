package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	auditmocks "github.com/Skywardkite/service-metrics/internal/audit/mocks"
	"github.com/Skywardkite/service-metrics/internal/config/server_config"
	"github.com/Skywardkite/service-metrics/internal/handler"
	model "github.com/Skywardkite/service-metrics/internal/model"
	servicemocks "github.com/Skywardkite/service-metrics/internal/service/mocks"
	pointers "github.com/Skywardkite/service-metrics/internal/utils"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func ExampleHandler_UpdateMetricsBatchJSONHandler() {
	// Создаём батч из двух метрик: Gauge и Counter
	batch := []model.Metrics{
		{ID: "requests", MType: model.Gauge, Value: pointers.To(float64(42.0))},
		{ID: "tasks_processed", MType: model.Counter, Delta: pointers.To(int64(7))},
	}

	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	serviceMock := servicemocks.NewMockMetricServiceInterface(ctrl)
	serviceMock.EXPECT().SaveMetricsBatch(gomock.Any(), batch).Return(nil)

	auditMock := auditmocks.NewMockAuditPublisherInterface(ctrl)
	auditMock.EXPECT().Publish(gomock.Any())

	h := handler.NewHandler(
		serviceMock,
		&server_config.Config{},
		&zap.SugaredLogger{},
		auditMock,
	)

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
