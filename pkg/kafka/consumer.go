package kafka

import (
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"fmt"
)

func Consumer() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
    	"bootstrap.servers":    "localhost:9092",
    	"group.id":             "myGroup",
    	"auto.offset.reset":    "smallest"})

	if err != nil {
        panic(fmt.Sprintf("Failed to create consumer: %v", err))
    }
   
    err = consumer.SubscribeTopics([]string{"orders_topic"}, nil)

    if err != nil {
            panic(err)
    }

    fmt.Println("Consumer initialized")
}
