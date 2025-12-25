package handler_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/Skywardkite/service-metrics/internal/config/server_config"
	"github.com/Skywardkite/service-metrics/internal/handler"
	mocks "github.com/Skywardkite/service-metrics/internal/handler/example_mocks"
	"go.uber.org/zap"
)

func ExampleHandler_GetMetric() {
	h := handler.NewHandler(&mocks.MockService{}, &server_config.Config{}, &mocks.MockStorage{}, &zap.SugaredLogger{}, &mocks.MockAudit{})

	req := httptest.NewRequest(http.MethodGet, "/value/gauge/requests", nil)
	w := httptest.NewRecorder()

	h.GetMetric(w, req)

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
