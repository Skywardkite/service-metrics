package handler_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	auditmocks "github.com/Skywardkite/service-metrics/internal/audit/mocks"
	"github.com/Skywardkite/service-metrics/internal/config/server_config"
	"github.com/Skywardkite/service-metrics/internal/handler"
	servicemocks "github.com/Skywardkite/service-metrics/internal/service/mocks"
	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func ExampleHandler_GetMetric() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	serviceMock := servicemocks.NewMockMetricServiceInterface(ctrl)
	serviceMock.EXPECT().GetMetric(gomock.Any(), "gauge", "requests").Return("42.0", nil)

	h := handler.NewHandler(
		serviceMock,
		&server_config.Config{},
		&zap.SugaredLogger{},
		&auditmocks.MockAuditPublisherInterface{},
	)

	r := chi.NewRouter()
	r.Get("/value/{metricType}/{metricName}", h.GetMetric)

	req := httptest.NewRequest(http.MethodGet, "/value/gauge/requests", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))
	fmt.Println("Body:", string(body))

	// Output:
	// Status Code: 200
	// Content-Type: application/json
	// Body: 42.0
}
