package handlers

import "net/http"

// Health ..
var Health = health{}

type health struct {
	handlers
}

// Get ...
func (h health) Get(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/health" {
		http.NotFound(w, r)
		return
	}
}
