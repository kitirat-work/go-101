//go:build integration
// +build integration

package dbtest

import (
	"context"
	"fmt"
	"time"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const defaultMySQLImage = "mysql:8.4" // LTS version ของ MySQL

type MySQLContainer struct {
	Container tc.Container

	Host     string
	Port     string // mapped port
	User     string
	Password string
	DB       string

	// DSN สำหรับ database/sql (go-sql-driver/mysql)
	// เช่น: user:pass@tcp(host:port)/db?parseTime=true
	DSN string
}

func (c *MySQLContainer) Stop(ctx context.Context) error {
	if c == nil || c.Container == nil {
		return nil
	}
	return c.Container.Terminate(ctx)
}

func StartMySQL(ctx context.Context, opts ...Option) (*MySQLContainer, error) {
	cfg := defaultConfigMySQL()
	for _, opt := range opts {
		opt(cfg)
	}

	env := map[string]string{
		"MYSQL_ROOT_PASSWORD": cfg.password,
		"MYSQL_USER":          cfg.user,
		"MYSQL_PASSWORD":      cfg.password,
		"MYSQL_DATABASE":      cfg.db,
	}
	for k, v := range cfg.env {
		env[k] = v
	}

	req := tc.ContainerRequest{
		Image:        cfg.image,
		Env:          env,
		ExposedPorts: []string{"3306/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("port: 3306  MySQL Community Server - GPL"),
			wait.ForListeningPort("3306/tcp"),
		).WithDeadline(cfg.waitDeadline),
		AlwaysPullImage: cfg.alwaysPull,
	}

	cont, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("start mysql container: %w", err)
	}

	host, err := cont.Host(ctx)
	if err != nil {
		_ = cont.Terminate(ctx)
		return nil, fmt.Errorf("get host: %w", err)
	}
	mp, err := cont.MappedPort(ctx, "3306/tcp")
	if err != nil {
		_ = cont.Terminate(ctx)
		return nil, fmt.Errorf("get mapped port: %w", err)
	}

	port := mp.Port()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.user, cfg.password, host, port, cfg.db)

	return &MySQLContainer{
		Container: cont,
		Host:      host,
		Port:      port,
		User:      cfg.user,
		Password:  cfg.password,
		DB:        cfg.db,
		DSN:       dsn,
	}, nil
}

// ------- Options -------

type Option func(*mysqlConfig)

type mysqlConfig struct {
	image        string
	user         string
	password     string
	db           string
	waitDeadline time.Duration
	alwaysPull   bool
	env          map[string]string
}

func defaultConfigMySQL() *mysqlConfig {
	return &mysqlConfig{
		image:        defaultMySQLImage,
		user:         "testuser",
		password:     "testpass",
		db:           "appdb",
		waitDeadline: 60 * time.Second,
		env:          map[string]string{},
	}
}

func WithMySQLImage(image string) Option {
	return func(c *mysqlConfig) { c.image = image }
}
func WithMySQLUser(user string) Option {
	return func(c *mysqlConfig) { c.user = user }
}
func WithMySQLPassword(pw string) Option {
	return func(c *mysqlConfig) { c.password = pw }
}
func WithMySQLDatabase(db string) Option {
	return func(c *mysqlConfig) { c.db = db }
}
func WithMySQLWaitDeadline(d time.Duration) Option {
	return func(c *mysqlConfig) { c.waitDeadline = d }
}
func WithMySQLAlwaysPull(always bool) Option {
	return func(c *mysqlConfig) { c.alwaysPull = always }
}
func WithMySQLEnv(k, v string) Option {
	return func(c *mysqlConfig) { c.env[k] = v }
}
