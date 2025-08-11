package main

import (
	"compress/gzip"
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
        
        // Создаем обертку для ResponseWriter
        writer := &gzipResponseWriter{
            ResponseWriter: w,
            request:        r,
        }
        defer writer.Close()

		next.ServeHTTP(writer, r)
	})
}

type gzipResponseWriter struct {
    http.ResponseWriter
    request     *http.Request
    gzipWriter  *gzip.Writer
    wroteHeader bool
}

func (w *gzipResponseWriter) WriteHeader(code int) {
    if w.wroteHeader {
        return
    }
    w.wroteHeader = true

    // Проверяем, нужно ли сжимать
    contentType := w.Header().Get("Content-Type")
    acceptsGzip := strings.Contains(w.request.Header.Get("Accept-Encoding"), "gzip")

    if acceptsGzip && (contentType == "text/html" || contentType == "application/json") {
        w.Header().Set("Content-Encoding", "gzip")
        w.gzipWriter = gzip.NewWriter(w.ResponseWriter)
    }

    w.ResponseWriter.WriteHeader(code)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
    if !w.wroteHeader {
        w.WriteHeader(http.StatusOK)
    }

    // Если Content-Type еще не установлен, пытаемся определить
    if w.Header().Get("Content-Type") == "" && len(b) > 0 {
        w.Header().Set("Content-Type", http.DetectContentType(b))
    }

    if w.gzipWriter != nil {
        return w.gzipWriter.Write(b)
    }
    return w.ResponseWriter.Write(b)
}

func (w *gzipResponseWriter) Close() {
    if w.gzipWriter != nil {
        w.gzipWriter.Close()
    }
}