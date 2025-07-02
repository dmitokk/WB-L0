package main

import (
	"log"
	"orders/pkg/cache"
	"orders/pkg/db"
	"orders/pkg/kafka"
	"orders/pkg/worker"
)

/*
1. Подключение к БД
2. Подключение к Кафке
3. Инициализация кэша
4. Восстановление заказов из БД в кэш
5. Запуск worker
6. Запуск http-сервера
*/

func main() {

    // 1. Подключение к БД
	database, err := db.Connect()
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }

    defer database.Close()

    // 2. Подключение к Кафке(инициализация консьюмера)
    kafka, err := kafka.ConsumerInit()
    if err != nil {
        log.Fatal("Failed to init kafka:", err)
    }

    defer kafka.Close()

    // 3. Инициализация кэша
    cache.Init()

    // 4. Восстановление заказов из БД в кэш при запуске
    err = cache.LoadFromDB(database)
    if err != nil {
        log.Fatal("Failed to create cache:", err)
    }

    // 5. Запуск worker процесса в фоне
    go worker.Start(kafka, database)
    
}