package db

import (
	"database/sql"
	"otp/internal/config"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySqlDB(cfg config.MySqlConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.Url)
	if err != nil {
		return nil, err
	}

	connMaxLifetime, err := time.ParseDuration(cfg.ConnMaxLifetime + "s")
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(connMaxLifetime)

	maxOpenConns, err := parseStringToInt(cfg.MaxOpenConns)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)

	maxIdleConns, err := parseStringToInt(cfg.MaxIdleConns)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(maxIdleConns)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func parseStringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}
