package http

import (
	"fmt"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 3 || parts[1] != "order" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	orderUID := parts[2]
	fmt.Println(orderUID)
	// get Order


	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}