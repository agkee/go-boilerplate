package web

import (
	// Add all /debug/pprof routes
	_ "net/http/pprof"

	"github.com/gorilla/mux"
)

// NewRouter instantiates a gorilla mux router and adds all relevant handlers
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ping", pingHandler)
	return r
}
