package handler

import (
	"net/http"
)

func (h *Handler) PingHandler(res http.ResponseWriter, req *http.Request) {
	if h.store == nil {
		h.logger.Errorw("Failed to chek connection store")
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := h.store.Ping(); err != nil {
		h.logger.Errorw("Failed to ping store", "error", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("OK"))
}