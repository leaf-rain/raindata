package rpprof

import (
	"net/http"
	"net/http/pprof"
)

// NewHandler new a pprof handler.
func NewPprofHandler(mux *http.ServeMux) {
	if mux == nil {
		mux = http.NewServeMux()
	}
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
