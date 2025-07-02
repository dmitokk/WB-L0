package http

import (
	"database/sql"
	"fmt"
	"net/http"
	"orders/pkg/cache"
	"orders/pkg/db"
	"strings"
	"encoding/json"
)

func handler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		parts := strings.Split(path, "/")

		if len(parts) < 3 || parts[1] != "order" {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		orderUID := parts[2]
		fmt.Println(orderUID)

		order, ok := cache.Get(orderUID)

		if !ok {
			orderPtr, err := db.GetOrderById(database, orderUID)
			if err != nil {

			}
			order = *orderPtr
		}
	
		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		data, err := json.Marshal(order)
    	if err != nil {
        	fmt.Println(err)
        	return
    	}

		w.Write(data)
	}
}