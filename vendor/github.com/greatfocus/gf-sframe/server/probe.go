package server

import (
	"net/http"
)

// liveProbe struct
type liveProbe struct{}

// ServeHTTP checks if is valid method
func (l liveProbe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		l.getProbe(w, r)
		return
	}
	// catch all
	// if no method is satisfied return an error
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Add("Allow", "GET")
}

// getFiles method
func (l *liveProbe) getProbe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
