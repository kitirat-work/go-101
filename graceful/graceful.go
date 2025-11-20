package graceful

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ILogger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

func newLogger() ILogger {
	return log.Default()
}

type GracefulShutdownConfig struct {
	ShutdownTimeoutSeconds int
	Signals                []os.Signal
	StartServer            func() error
	StopServer             func(ctx context.Context) error
	Logger                 ILogger
}

func Run(cfg GracefulShutdownConfig) {
	logger := cfg.Logger
	if logger == nil {
		logger = newLogger()
	}
	if cfg.StartServer == nil {
		logger.Println("[graceful] StartServer is nil")
		return
	}
	if len(cfg.Signals) == 0 {
		cfg.Signals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
	}

	// 1) ดักสัญญาณ
	ctx, stop := signal.NotifyContext(context.Background(), cfg.Signals...)
	defer stop()

	// 2) รันเซิร์ฟเวอร์ใน goroutine
	srvErr := make(chan error, 1)
	go func() {
		err := cfg.StartServer()
		// ถ้าเป็น net/http แนะนำเมิน http.ErrServerClosed
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErr <- err
		} else {
			srvErr <- nil
		}
	}()

	// 3) รอ: สัญญาณ หรือ เซิร์ฟเวอร์พังเอง
	select {
	case <-ctx.Done():
		logger.Println("[graceful] signal received → start shutdown")
	case err := <-srvErr:
		if err != nil {
			logger.Printf("[graceful] server exited with error: %v\n", err)
		} else {
			logger.Println("[graceful] server exited normally")
		}
		return
	}

	// 4) เมื่อเริ่ม shutdown: ตั้ง timer context
	timeout := time.Duration(cfg.ShutdownTimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		// fallback: เรียก StopServer แบบเดิม
		if cfg.StopServer != nil {
			if err := cfg.StopServer(shutdownCtx); err != nil {
				logger.Printf("[graceful] stop error: %v\n", err)
			}
		}
	}()

	select {
	case <-done:
		logger.Println("[graceful] shutdown completed")
	case <-shutdownCtx.Done():
		logger.Printf("[graceful] timeout (%s) exceeded → force exit\n", timeout)
	}
}
