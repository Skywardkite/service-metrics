package handler

import (
	"net/http"
)

func (h *Handler) PingHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if err := h.service.Ping(ctx); err != nil {
		h.logger.Errorw("Failed to ping store", "error", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responseBody := []byte("OK")
	if h.cfg.Key != "" {
		hash := SignBody(responseBody, h.cfg.Key)
		res.Header().Set("HashSHA256", hash)
	}
	res.WriteHeader(http.StatusOK)
	res.Write(responseBody)
}
