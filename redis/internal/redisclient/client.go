package redisclient

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	once   sync.Once
	client *redis.Client
)

// New returns a process-wide singleton client (safe for concurrent use).
// go-redis manages connection pooling internally, so reuse a single client.
func New(cfg Config) *redis.Client {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:            cfg.Addr,
			Password:        cfg.Password,
			DB:              cfg.DB,
			PoolSize:        cfg.PoolSize,
			MinIdleConns:    cfg.MinIdleConns,
			DialTimeout:     cfg.DialTimeout,
			ReadTimeout:     cfg.ReadTimeout,
			WriteTimeout:    cfg.WriteTimeout,
			ConnMaxIdleTime: cfg.ConnMaxIdleTime,
			ConnMaxLifetime: cfg.ConnMaxLifetime,
			// Protocol: 2, // ตั้งได้หากต้องการบังคับ ตามตัวอย่างหน้า Connect (ค่า default จะจัดการอัตโนมัติ)
		})
	})
	return client
}

// PingWithTimeout ตรวจสุขภาพการเชื่อมต่อ
func PingWithTimeout(ctx context.Context, c *redis.Client, t time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, t)
	defer cancel()
	return c.Ping(ctx).Err()
}

// Close ใช้ตอน graceful shutdown
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
