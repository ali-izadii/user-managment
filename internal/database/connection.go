package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

var Conn *pgx.Conn

func InitPostgres(uri string) {
	var err error
	Conn, err = pgx.Connect(context.Background(), uri)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Test the connection
	var version string
	if err := Conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	log.Println("Connected to:", version)
}

func Close() {
	if Conn != nil {
		Conn.Close(context.Background())
	}
}
