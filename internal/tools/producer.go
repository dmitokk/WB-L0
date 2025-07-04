package main

import (
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
    	"bootstrap.servers":	"localhost:9092",
    	"acks":             	"all",
    	"client.id":    		"myProducer"})

	if err != nil {
    	panic(err)
	}

	defer producer.Close()
	fmt.Println("Producer initialized")

	dirPath := "examples/orders_json"

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	var jsonData [][]byte
	for _, entry := range entries {
		if !entry.IsDir() {
			data, err := os.ReadFile(filepath.Join(dirPath, entry.Name()))
			if err != nil {
				panic(err)
			}
			jsonData = append(jsonData, data)
		}
	}

	go func() {
        for e := range producer.Events() {
            switch ev := e.(type) {
            case *kafka.Message:
                if ev.TopicPartition.Error != nil {
                    fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition.Error)
                } else {
                	fmt.Printf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
                    *ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
                }
            }
        }
    }()

	topic := "orders_topic"
	for _, json := range jsonData {
		producer.Produce(&kafka.Message{
            TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
            Value:          []byte(json),
        }, nil)
	}

}
