package redisclient

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr            string // host:port
	Password        string
	DB              int
	PoolSize        int
	MinIdleConns    int
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseInt(key string, def int) int {
	if s := os.Getenv(key); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			return n
		}
	}
	return def
}

func parseDuration(key string, def time.Duration) time.Duration {
	if s := os.Getenv(key); s != "" {
		if d, err := time.ParseDuration(s); err == nil {
			return d
		}
	}
	return def
}

func LoadConfig() Config {
	return Config{
		Addr:            getenv("REDIS_ADDR", "127.0.0.1:6379"),
		Password:        getenv("REDIS_PASSWORD", ""),
		DB:              parseInt("REDIS_DB", 0),
		PoolSize:        parseInt("REDIS_POOL_SIZE", 20),
		MinIdleConns:    parseInt("REDIS_MIN_IDLE", 5),
		DialTimeout:     parseDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		ReadTimeout:     parseDuration("REDIS_READ_TIMEOUT", 2*time.Second),
		WriteTimeout:    parseDuration("REDIS_WRITE_TIMEOUT", 2*time.Second),
		ConnMaxIdleTime: parseDuration("REDIS_CONN_MAX_IDLE", 30*time.Minute),
		ConnMaxLifetime: parseDuration("REDIS_CONN_MAX_LIFETIME", 2*time.Hour),
	}
}
