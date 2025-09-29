//go:build integration
// +build integration

package dbtest

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CreateTablesFromModels(dsn string, tableModels ...any) error {
	gormDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect to db: %w", err)
	}
	defer func() {
		sqlDB, err := gormDb.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	err = gormDb.AutoMigrate(tableModels...)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}
