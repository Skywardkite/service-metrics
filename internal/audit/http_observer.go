package audit

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type HttpObserver struct {
	URL string
}

func (h *HttpObserver) Notify(event AuditEvent) {
	data, _ := json.Marshal(event)
	http.Post(h.URL, "application/json", bytes.NewReader(data))
}
