package main

import (
	"log"
	"orders/pkg/db"
    "orders/pkg/kafka"
)

func main() {

    // 1. Connecting to Postgres
	db, err := db.Connect()
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }

    defer db.Close()

    // 2. Initialising Kafka consumer
    kafka, err := kafka.ConsumerInit()
    if err != nil {
        log.Fatal("Failed to init kafka:", err)
    }

    defer kafka.Close()

}