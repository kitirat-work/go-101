//go:build integration
// +build integration

package db_test

import (
	"context"
	"testing"
	"time"

	"otp/internal/config"
	"otp/internal/db"
	"otp/internal/dbtest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMySQLDBSuccess(t *testing.T) {
	// arrange
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	mc, err := dbtest.StartMySQL(ctx,
		dbtest.WithMySQLDatabase("testdb"),
		dbtest.WithMySQLUser("testuser"),
		dbtest.WithMySQLPassword("testpass"),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = mc.Stop(context.Background()) })
	cfg := config.MySqlConfig{
		Url:             mc.DSN,
		ConnMaxLifetime: "1", // 1s
		MaxOpenConns:    "5",
		MaxIdleConns:    "2",
	}

	// act 1
	dbConn, err := db.NewMySqlDB(cfg)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, dbConn)
	t.Cleanup(func() { dbConn.Close() })

	// act 2 Check connection is alive
	err = dbConn.PingContext(ctx)
	assert.NoError(t, err)

	// act 3 Simple query
	_, err = dbConn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS t (id INT PRIMARY KEY AUTO_INCREMENT, v INT NOT NULL)`)
	assert.NoError(t, err)
}

func TestNewMySQLDBInvalidDSN(t *testing.T) {
	cfg := config.MySqlConfig{
		Url:             "invalid-dsn",
		ConnMaxLifetime: "1",
		MaxOpenConns:    "1",
		MaxIdleConns:    "1",
	}
	dbConn, err := db.NewMySqlDB(cfg)
	require.Error(t, err)
	require.Nil(t, dbConn)
}

func TestNewMySQLDBInvalidConns(t *testing.T) {
	const notANumber = "not-a-number"
	cfg := config.MySqlConfig{
		Url:             "user:pass@tcp(localhost:3306)/db",
		ConnMaxLifetime: notANumber,
		MaxOpenConns:    notANumber,
		MaxIdleConns:    notANumber,
	}
	dbConn, err := db.NewMySqlDB(cfg)
	require.Error(t, err)
	require.Nil(t, dbConn)
}
