package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func gzipMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		contentEncoding := r.Header.Get("Content-Encoding")
		isGzipped := strings.Contains(contentEncoding, "gzip")

		if isGzipped {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gz.Close()
			r.Body = gz
		}

		var gzWriter *gzip.Writer
		if supportsGzip {
			contentType := r.Header.Get("Content-Type")
			if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html") {
				w.Header().Set("Content-Encoding", "gzip")
				gzWriter = gzip.NewWriter(w)
				defer gzWriter.Close()
				w = &gzipResponseWriter{Writer: gzWriter, ResponseWriter: w}
			}
		}

		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}