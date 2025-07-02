package worker

import (
    "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"orders/pkg/model"
	"encoding/json"
	"orders/pkg/db"
	"database/sql"
	"orders/pkg/cache"
	"fmt"
	"log"
)

/*
Основной рабочий цикл:

1. Получение сообщения
2. Парсинг
3. Сохранение в БД
4. Обновление кеша
5. Подтверждение коммита(смещения)
*/

func Start(consumer *kafka.Consumer, database *sql.DB) error {

	for {
		ev := consumer.Poll(100)
		switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
				var order model.Order
				err := json.Unmarshal(e.Value, &order)
    			if err != nil {
					log.Printf("failed to unmarshal order: %s\nraw: %s\n", err, string(e.Value))
					return fmt.Errorf("error unable to marshal JSON due to : %w", err)
    			}
				if order.OrderUID == "" || order.Delivery == nil || order.Payment == nil {
					return fmt.Errorf("invalid order: missing required fields")
				}
				err = db.InsertOrder(database, order)

				if err != nil {
					return fmt.Errorf("invalid order: missing required fields")
				}

				cache.Set(order)

				consumer.Commit()
			case kafka.Error:
				return fmt.Errorf("error while reading messages in consumer: %w", e)
			default:
				fmt.Printf("Ignored %v\n", e)
		}
	}
}