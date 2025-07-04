package cache

import (
	"database/sql"
	"orders/pkg/model"
	"sync"
	"orders/pkg/db"
	"fmt"
)

type Cache struct {
  mu sync.RWMutex
  data map[string]model.Order
}

var c Cache

func Init() {
	c = Cache{
		data: make(map[string]model.Order),
	}
	fmt.Println("Cache initialized!")
}

func Set(order model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[order.OrderUID] = order
}

func Get(orderUID string) (model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.data[orderUID]
	return order, ok
}

func LoadFromDB(database *sql.DB) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	orders, err := db.GetAllOrders(database)
	if err != nil {
		return fmt.Errorf("error getting all orders: %w", err)
	}

	for _, order := range orders {
		c.data[order.OrderUID] = *order
	}

	fmt.Println("Cache loaded from db!")
	return nil
}