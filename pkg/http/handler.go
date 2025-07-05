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
		fmt.Print(order)
		if !ok {
			orderPtr, err := db.GetOrderById(database, orderUID)
			if err != nil {
				fmt.Print("Thiss order does not exist!")
			}
			order = *orderPtr
		}
	
		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)

		data, err := json.Marshal(order)
    	if err != nil {
        	fmt.Println(err)
        	return
    	}

		w.Write(data)
	}
}