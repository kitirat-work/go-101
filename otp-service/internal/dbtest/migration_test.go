//go:build integration
// +build integration

package dbtest

import (
	"context"
	"otp/internal/config"
	"otp/internal/db"
	"otp/internal/entities"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateTablesFromModels(t *testing.T) {
	t.Run("should create tables successfully", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		// Start MySQL container
		mc, err := StartMySQL(ctx, WithMySQLDatabase("migrtest"), WithMySQLUser("miguser"), WithMySQLPassword("migpass"))
		assert.NoError(t, err)
		defer mc.Stop(context.Background())

		// Run migrations
		err = CreateTablesFromModels(mc.DSN, &entities.OtpCode{}, &entities.Session{})

		assert.NoError(t, err)
		assert.NotNil(t, mc)

		// Verify tables exist
		sqlDb, err := db.NewMySqlDB(config.MySqlConfig{
			Url:             mc.DSN,
			ConnMaxLifetime: "1",
			MaxOpenConns:    "1",
			MaxIdleConns:    "1",
		})
		assert.NoError(t, err)
		defer sqlDb.Close()

		// count rows in otp_code table
		var count int
		err = sqlDb.QueryRowContext(ctx, "SELECT COUNT(*) FROM otp_code").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)

		// count rows in session table
		err = sqlDb.QueryRowContext(ctx, "SELECT COUNT(*) FROM session").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("should fail with invalid DSN", func(t *testing.T) {
		err := CreateTablesFromModels("invalid-dsn", &entities.OtpCode{})
		assert.Error(t, err)
	})
}
