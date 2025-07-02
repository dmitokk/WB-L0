package http

import (
	"fmt"
	"net/http"
	"database/sql"
)

func Server(database *sql.DB) error {
	http.HandleFunc("/order/", handler(database))

	if err := http.ListenAndServe(":8081", nil); err != nil {
        return fmt.Errorf("error starting server: %w", err)
    }

	return nil
}