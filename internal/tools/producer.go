package tools

import (
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"fmt"
)

func Producer() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
    	"bootstrap.servers":	"localhost:9092",
    	"acks":             	"all",
    	"client.id":    		"myProducer"})

	if err != nil {
    	panic(err)
	}

	defer producer.Close()
	fmt.Println("Producer initialized")

}
