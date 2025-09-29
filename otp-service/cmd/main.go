package main

import (
	"log"
	"otp/internal/config"
	"otp/internal/db"
)

func main() {
	cfg := config.NewConfig()

	dbClient, err := db.NewMySqlDB(cfg.MySqlConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbClient.Close()

}
