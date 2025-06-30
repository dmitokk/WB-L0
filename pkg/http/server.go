package http

import (
	"fmt"
	"net/http"
)

func Server() error {
	http.HandleFunc("/order/", handler)

	if err := http.ListenAndServe(":8081", nil); err != nil {
        return fmt.Errorf("error starting server: %w", err)
    }

	return nil
}