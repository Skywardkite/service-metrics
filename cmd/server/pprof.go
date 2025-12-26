package server

import (
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
)

// MountPprof создает маршруты профилировщика, чтобы смотреть pprof-профили в браузере.
func MountPprof(r chi.Router) {
	r.HandleFunc("/", pprof.Index)
	r.HandleFunc("/cmdline", pprof.Cmdline)
	r.HandleFunc("/profile", pprof.Profile)
	r.HandleFunc("/symbol", pprof.Symbol)
	r.HandleFunc("/trace", pprof.Trace)
	r.HandleFunc("/allocs", pprof.Handler("allocs").ServeHTTP)
	r.HandleFunc("/block", pprof.Handler("block").ServeHTTP)
	r.HandleFunc("/goroutine", pprof.Handler("goroutine").ServeHTTP)
	r.HandleFunc("/heap", pprof.Handler("heap").ServeHTTP)
	r.HandleFunc("/mutex", pprof.Handler("mutex").ServeHTTP)
	r.HandleFunc("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}
