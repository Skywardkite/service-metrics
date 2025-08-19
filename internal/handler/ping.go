package handler

import (
	"net/http"
)

func (h *Handler) PingHandler(res http.ResponseWriter, req *http.Request) {
	if h.db == nil {
		http.Error(res, "Database not configured", http.StatusInternalServerError)
		return
	}

	if err := h.db.Ping(); err != nil {
		http.Error(res, "Database connection failed", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("OK"))
}