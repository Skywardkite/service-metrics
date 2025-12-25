package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/Skywardkite/service-metrics/internal/config/server_config"
	"github.com/Skywardkite/service-metrics/internal/handler"
	mocks "github.com/Skywardkite/service-metrics/internal/handler/example_mocks"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func ExampleHandler_UpdateHandler() {
	h := handler.NewHandler(&mocks.MockService{}, &server_config.Config{}, &mocks.MockStorage{}, &zap.SugaredLogger{}, &mocks.MockAudit{})

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
