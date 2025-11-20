package graceful_test

import (
	"context"
	"errors"
	"fmt"
	"graceful"

	"net/http"
	"os"

	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock logger to capture log output
type mockLogger struct {
	mu      sync.Mutex
	entries []string
}

func (m *mockLogger) Println(v ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, strings.TrimSpace(fmt.Sprint(v...)))
}

func (m *mockLogger) Printf(format string, v ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, strings.TrimSpace(fmt.Sprintf(format, v...)))
}

func (m *mockLogger) getEntries() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	entries := make([]string, len(m.entries))
	copy(entries, m.entries)
	return entries
}

func TestRun_SuccessfulServerStartAndGracefulShutdown(t *testing.T) {
	logger := &mockLogger{}
	serverStarted := make(chan struct{})
	serverStopped := make(chan struct{})
	shutdownCalled := make(chan struct{})

	cfg := graceful.GracefulShutdownConfig{
		ShutdownTimeoutSeconds: 5,
		StartServer: func() error {
			close(serverStarted)
			<-serverStopped // Block until we signal to stop
			return http.ErrServerClosed
		},
		StopServer: func(ctx context.Context) error {
			close(shutdownCalled)
			close(serverStopped)
			return nil
		},
		Logger: logger,
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		graceful.Run(cfg)
	}()

	// Wait for server to start
	select {
	case <-serverStarted:
	case <-time.After(time.Second):
		t.Fatal("server should have started")
	}

	// Send interrupt signal
	go func() {
		time.Sleep(100 * time.Millisecond)
		process, _ := os.FindProcess(os.Getpid())
		process.Signal(os.Interrupt)
	}()

	// Wait for shutdown to be called
	select {
	case <-shutdownCalled:
	case <-time.After(2 * time.Second):
		t.Fatal("shutdown should have been called")
	}

	// Wait for Run to complete
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run should have completed")
	}

	entries := logger.getEntries()
	assert.Contains(t, strings.Join(entries, " "), "signal received")
	assert.Contains(t, strings.Join(entries, " "), "shutdown completed")
}

func TestRun_ServerExitsNormally(t *testing.T) {
	logger := &mockLogger{}
	serverStarted := make(chan struct{})

	cfg := graceful.GracefulShutdownConfig{
		StartServer: func() error {
			close(serverStarted)
			return http.ErrServerClosed // Normal HTTP server shutdown
		},
		Logger: logger,
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		graceful.Run(cfg)
	}()

	// Wait for server to start and exit
	select {
	case <-serverStarted:
	case <-time.After(time.Second):
		t.Fatal("server should have started")
	}

	// Wait for Run to complete
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run should have completed")
	}

	entries := logger.getEntries()
	assert.Contains(t, strings.Join(entries, " "), "server exited normally")
}

func TestRun_ServerExitsWithError(t *testing.T) {
	logger := &mockLogger{}
	expectedError := errors.New("server startup failed")

	cfg := graceful.GracefulShutdownConfig{
		StartServer: func() error {
			return expectedError
		},
		Logger: logger,
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		graceful.Run(cfg)
	}()

	// Wait for Run to complete
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run should have completed")
	}

	entries := logger.getEntries()
	assert.Contains(t, strings.Join(entries, " "), "server exited with error")
	assert.Contains(t, strings.Join(entries, " "), expectedError.Error())
}

func TestRun_NilStartServer(t *testing.T) {
	logger := &mockLogger{}

	cfg := graceful.GracefulShutdownConfig{
		StartServer: nil,
		Logger:      logger,
	}

	graceful.Run(cfg)

	entries := logger.getEntries()
	assert.Contains(t, strings.Join(entries, " "), "StartServer is nil")
}

func TestRun_DefaultTimeout(t *testing.T) {
	logger := &mockLogger{}
	serverStarted := make(chan struct{})
	serverStopped := make(chan struct{})

	cfg := graceful.GracefulShutdownConfig{
		ShutdownTimeoutSeconds: 0, // Should default to 30 seconds
		StartServer: func() error {
			close(serverStarted)
			<-serverStopped
			return http.ErrServerClosed
		},
		StopServer: func(ctx context.Context) error {
			// Check that timeout is 30 seconds
			deadline, ok := ctx.Deadline()
			if ok {
				timeout := time.Until(deadline)
				assert.True(t, timeout > 25*time.Second && timeout <= 30*time.Second)
			}
			close(serverStopped)
			return nil
		},
		Logger: logger,
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		graceful.Run(cfg)
	}()

	// Wait for server to start
	select {
	case <-serverStarted:
	case <-time.After(time.Second):
		t.Fatal("server should have started")
	}

	// Send signal
	go func() {
		time.Sleep(100 * time.Millisecond)
		process, _ := os.FindProcess(os.Getpid())
		process.Signal(syscall.SIGTERM)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run should have completed")
	}
}

func TestRun_StopServerError(t *testing.T) {
	logger := &mockLogger{}
	serverStarted := make(chan struct{})
	serverStopped := make(chan struct{})
	stopError := errors.New("stop server failed")

	cfg := graceful.GracefulShutdownConfig{
		ShutdownTimeoutSeconds: 2,
		StartServer: func() error {
			close(serverStarted)
			<-serverStopped
			return http.ErrServerClosed
		},
		StopServer: func(ctx context.Context) error {
			close(serverStopped)
			return stopError
		},
		Logger: logger,
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		graceful.Run(cfg)
	}()

	// Wait for server to start
	select {
	case <-serverStarted:
	case <-time.After(time.Second):
		t.Fatal("server should have started")
	}

	// Send signal
	go func() {
		time.Sleep(100 * time.Millisecond)
		process, _ := os.FindProcess(os.Getpid())
		process.Signal(os.Interrupt)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Run should have completed")
	}

	entries := logger.getEntries()
	assert.Contains(t, strings.Join(entries, " "), "stop error")
	assert.Contains(t, strings.Join(entries, " "), stopError.Error())
}

func TestRun_ShutdownTimeout(t *testing.T) {
	logger := &mockLogger{}
	serverStarted := make(chan struct{})
	serverStopped := make(chan struct{})

	cfg := graceful.GracefulShutdownConfig{
		ShutdownTimeoutSeconds: 1, // Short timeout
		StartServer: func() error {
			close(serverStarted)
			<-serverStopped
			return http.ErrServerClosed
		},
		StopServer: func(ctx context.Context) error {
			// Don't close serverStopped to simulate hanging shutdown
			time.Sleep(2 * time.Second) // Longer than timeout
			close(serverStopped)
			return nil
		},
		Logger: logger,
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		graceful.Run(cfg)
	}()

	// Wait for server to start
	select {
	case <-serverStarted:
	case <-time.After(time.Second):
		t.Fatal("server should have started")
	}

	// Send signal
	go func() {
		time.Sleep(100 * time.Millisecond)
		process, _ := os.FindProcess(os.Getpid())
		process.Signal(os.Interrupt)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Run should have completed")
	}

	entries := logger.getEntries()
	assert.Contains(t, strings.Join(entries, " "), "timeout")
	assert.Contains(t, strings.Join(entries, " "), "force exit")
}

func TestRun_NilStopServer(t *testing.T) {
	logger := &mockLogger{}
	serverStarted := make(chan struct{})
	serverStopped := make(chan struct{})

	cfg := graceful.GracefulShutdownConfig{
		ShutdownTimeoutSeconds: 2,
		StartServer: func() error {
			close(serverStarted)
			<-serverStopped
			return http.ErrServerClosed
		},
		StopServer: nil, // No stop server function
		Logger:     logger,
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		graceful.Run(cfg)
	}()

	// Wait for server to start
	select {
	case <-serverStarted:
	case <-time.After(time.Second):
		t.Fatal("server should have started")
	}

	// Send signal and stop server manually since StopServer is nil
	go func() {
		time.Sleep(100 * time.Millisecond)
		process, _ := os.FindProcess(os.Getpid())
		process.Signal(os.Interrupt)
		// Manually stop the server since StopServer is nil
		time.Sleep(200 * time.Millisecond)
		close(serverStopped)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Run should have completed")
	}

	entries := logger.getEntries()
	assert.Contains(t, strings.Join(entries, " "), "shutdown completed")
}

func TestRun_CustomSignals(t *testing.T) {
	logger := &mockLogger{}
	serverStarted := make(chan struct{})
	serverStopped := make(chan struct{})

	cfg := graceful.GracefulShutdownConfig{
		ShutdownTimeoutSeconds: 2,
		Signals:                []os.Signal{syscall.SIGUSR1}, // Custom signal
		StartServer: func() error {
			close(serverStarted)
			<-serverStopped
			return http.ErrServerClosed
		},
		StopServer: func(ctx context.Context) error {
			close(serverStopped)
			return nil
		},
		Logger: logger,
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		graceful.Run(cfg)
	}()

	// Wait for server to start
	select {
	case <-serverStarted:
	case <-time.After(time.Second):
		t.Fatal("server should have started")
	}

	// Send custom signal
	go func() {
		time.Sleep(100 * time.Millisecond)
		process, _ := os.FindProcess(os.Getpid())
		process.Signal(syscall.SIGUSR1)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Run should have completed")
	}

	entries := logger.getEntries()
	assert.Contains(t, strings.Join(entries, " "), "signal received")
	assert.Contains(t, strings.Join(entries, " "), "shutdown completed")
}

func TestRun_NegativeTimeout(t *testing.T) {
	logger := &mockLogger{}
	serverStarted := make(chan struct{})
	serverStopped := make(chan struct{})

	cfg := graceful.GracefulShutdownConfig{
		ShutdownTimeoutSeconds: -5, // Negative timeout should default to 30s
		StartServer: func() error {
			close(serverStarted)
			<-serverStopped
			return http.ErrServerClosed
		},
		StopServer: func(ctx context.Context) error {
			// Verify timeout is 30 seconds
			deadline, ok := ctx.Deadline()
			if ok {
				timeout := time.Until(deadline)
				assert.True(t, timeout > 25*time.Second && timeout <= 30*time.Second)
			}
			close(serverStopped)
			return nil
		},
		Logger: logger,
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		graceful.Run(cfg)
	}()

	// Wait for server to start
	select {
	case <-serverStarted:
	case <-time.After(time.Second):
		t.Fatal("server should have started")
	}

	// Send signal
	go func() {
		time.Sleep(100 * time.Millisecond)
		process, _ := os.FindProcess(os.Getpid())
		process.Signal(os.Interrupt)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Run should have completed")
	}
}
