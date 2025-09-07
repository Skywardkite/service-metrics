package handler

import (
	"net/http"
)

func (h *Handler) PingHandler(res http.ResponseWriter, req *http.Request) {
	if h.service.Cfg.Key != "" {
		success, err := checkKey(req, h.service.Cfg.Key)
		if err != nil {
			h.logger.Errorf("Error read body", err)
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if !success {
			h.logger.Info("no rights")
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

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

	responseBody := []byte("OK")

	if h.service.Cfg.Key != "" {
        hash := signBody(responseBody, h.service.Cfg.Key)
        res.Header().Set("HashSHA256", hash)
    }

	res.WriteHeader(http.StatusOK)
	res.Write(responseBody)
}