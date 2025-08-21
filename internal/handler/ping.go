package handler

import (
	"net/http"
)

func (h *Handler) PingHandler(res http.ResponseWriter, req *http.Request) {
	if h.store == nil {
		http.Error(res, "Database not configured", http.StatusInternalServerError)
		return
	}

	if err := h.store.Ping(); err != nil {
		http.Error(res, "Database connection failed", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("OK"))
}