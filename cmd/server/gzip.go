package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func gzipMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		 // Обработка входящего сжатого контента
        if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
            gz, err := gzip.NewReader(r.Body)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            defer gz.Close()
            r.Body = gz
        }
        
        // Перехватыватываем Content-Type перед отправкой заголовков
        acceptsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

		writer := &responseSniffer{ResponseWriter: w, acceptsGzip: acceptsGzip}

		next.ServeHTTP(writer, r)
	})
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}
type responseSniffer struct {
	http.ResponseWriter
	acceptsGzip bool
	contentType string
}

func (r *responseSniffer) WriteHeader(status int) {
	if r.contentType == "" {
		r.contentType = r.Header().Get("Content-Type")
	}
	
	if r.acceptsGzip && (r.contentType == "text/html" || r.contentType == "application/json") {
		r.Header().Set("Content-Encoding", "gzip")
		r.ResponseWriter = &gzipResponseWriter{
			Writer:         gzip.NewWriter(r.ResponseWriter),
			ResponseWriter: r.ResponseWriter,
		}
	}
	
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseSniffer) Write(b []byte) (int, error) {
	// Если Content-Type еще не установлен, пытаемся определить его
	if r.contentType == "" && len(b) > 0 {
		r.contentType = http.DetectContentType(b)
	}
	return r.ResponseWriter.Write(b)
}