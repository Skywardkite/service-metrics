package handler

import (
	"net"
	"net/http"
	"strings"
)

// clientIP узнает ip клиента из заголовка запроса для отправки аудита.
func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return strings.Split(fwd, ",")[0]
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}
