//go:build integration
// +build integration

package dbtest

import (
	"context"
	"database/sql"
	"otp/internal/config"
	"otp/internal/db"
	"otp/internal/entities"
	"time"
)

func CreateTestDB() (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Start MySQL container
	mc, err := StartMySQL(ctx, WithMySQLDatabase("testdb"), WithMySQLUser("test"), WithMySQLPassword("test"))
	if err != nil {
		return nil, err
	}

	// Run migrations
	err = CreateTablesFromModels(
		mc.DSN,
		&entities.Session{},
		&entities.OtpCode{},
	)
	if err != nil {
		_ = mc.Stop(context.Background())
		return nil, err
	}

	// Create DB connection
	cfg := config.MySqlConfig{
		Url:             mc.DSN,
		ConnMaxLifetime: "1", // 1s
		MaxOpenConns:    "5",
		MaxIdleConns:    "2",
	}
	dbConn, err := db.NewMySqlDB(cfg)
	if err != nil {
		_ = mc.Stop(context.Background())
		return nil, err
	}

	return dbConn, nil
}
