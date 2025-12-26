package handler_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func ExampleHandler_GetAllMetricsHandler() {
	// Хендлер с фиктивными данными
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gauges := map[string]float64{"requests": 42.0}
		counters := map[string]int64{"tasks_processed": 7}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Gauges: %v, Counters: %v", gauges, counters)
	})

	// Имитация HTTP запроса
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	testHandler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))
	fmt.Println("Body:", string(body))

	// Output:
	// Status Code: 200
	// Content-Type: text/html
	// Body: Gauges: map[requests:42], Counters: map[tasks_processed:7]
}
