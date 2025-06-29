package main

import (
	"log"
	"orders/pkg/db"
)

func main() {
	db, err := db.Connect()
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }

    defer db.Close()
}