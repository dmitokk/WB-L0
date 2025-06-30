package kafka

import (
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"fmt"
)

func ConsumerInit() (*kafka.Consumer, error){
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
    	"bootstrap.servers":    "localhost:9092",
    	"group.id":             "myGroup",
    	"auto.offset.reset":    "smallest"})

	if err != nil {
        return nil, fmt.Errorf("failed to create consumer: %w", err)
    }
   
    err = consumer.SubscribeTopics([]string{"orders_topic"}, nil)

    if err != nil {
        return nil, fmt.Errorf("error subscribing to topic: %w", err)
    }

    fmt.Println("Consumer initialized")

    return consumer, nil
}
