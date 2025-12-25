package audit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Skywardkite/service-metrics/internal/logger"
)

type HttpObserver struct {
	URL string
	mu  sync.Mutex
}

func (h *HttpObserver) Notify(event AuditEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		logger.Sugar.Errorw("Can't marshal audit event", "error", err)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	resp, err := http.Post(h.URL, "application/json", bytes.NewReader(data))
	if err != nil {
		logger.Sugar.Errorw("Can't send audit event", "error", err)
		return
	}

	defer resp.Body.Close()
}
