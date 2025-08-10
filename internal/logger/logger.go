package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var Sugar *zap.SugaredLogger

type responseWriterWrapper struct {
	http.ResponseWriter
	status int
	size   int
}

func Initialize() error {
    logger, err := zap.NewDevelopment()
    if err != nil {
        return err
    }

    Sugar = logger.Sugar()
	return nil
}

func Sync() error {
	return Sugar.Sync()
}

// WithLogging middleware для логирования запросов и ответов
func WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &responseWriterWrapper{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		h.ServeHTTP(ww, r)

		duration := time.Since(start)

		Sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"duration", duration,
			"status", ww.status,
			"size", ww.size,
		)
	})
}

func (w *responseWriterWrapper) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}