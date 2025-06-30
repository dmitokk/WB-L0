package db

import (
    "database/sql"
    "fmt"

    _ "github.com/lib/pq"
)

const (
    host     = "localhost"
    port     = 25432
    user     = "demo_user"
    password = "password"
    dbname   = "demo_orders"
)

func Connect() (*sql.DB, error) {
    // Establish a connection to the PostgreSQL database
    db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	
    if err != nil {
		  return nil, fmt.Errorf("error opening database: %w", err)
    }

    // Check if the connection is successful
    if err := db.Ping(); err != nil {
      db.Close() // clean up before returning
      return nil, fmt.Errorf("error pinging database: %w", err)
    }

    fmt.Println("Successfully connected to PostgreSQL!")
	
	return db, nil
}