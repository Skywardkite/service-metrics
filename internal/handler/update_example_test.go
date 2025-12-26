package handler_test

import (
	"fmt"
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

func ExampleHandler_UpdateHandler() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	serviceMock := servicemocks.NewMockMetricServiceInterface(ctrl)
	serviceMock.EXPECT().UpdateMetric(gomock.Any(), "gauge", "requests", "42").Return(nil)

	auditMock := auditmocks.NewMockAuditPublisherInterface(ctrl)
	auditMock.EXPECT().Publish(gomock.Any())

	//repMock := repmocks.NewMockStorage(nil)

	h := handler.NewHandler(
		serviceMock,
		&server_config.Config{},
		&zap.SugaredLogger{},
		auditMock,
	)

	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/update/gauge/requests/42", "text/plain", nil)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)

	// Output:
	// Status Code: 200
}
