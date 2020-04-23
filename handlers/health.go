package handlers

import "net/http"

// Health ..
var Health = health{}

type health struct {
	handlers
}

// Get ...
func (h health) Get(w http.ResponseWriter, r *http.Request) {

}
