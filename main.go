package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
	"vibrain/internal/core/queue"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/logger"
	"vibrain/internal/port/httpserver"

	"github.com/labstack/echo/v4"
)

func newQueue(ctx context.Context) *queue.Queue {
	q, err := queue.New(config.Settings.QueueDatabaseURL)
	if err != nil {
		logger.Default.Fatal("failed to create queue", "error", err)
	}
	go func() {
		if err := q.Start(ctx); err != nil {
			logger.Default.Fatal("failed to start queue", "error", err)
		}
		logger.Default.Info("queue started")
	}()
	return q
}

func newServer() *echo.Echo {
	server := httpserver.New()
	go func() {
		addr := fmt.Sprintf(":%d", config.Settings.Port)
		if err := server.Start(fmt.Sprintf(":%d", config.Settings.Port)); err != nil {
			logger.Default.Fatal("failed to start server", "addr", addr, "error", err)
		}
		logger.Default.Info("server started", "addr", addr)
	}()
	return server
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	logger.Default.Info("starting service")

	// start queue
	q := newQueue(ctx)
	// start http server
	server := newServer()

	// wait for signal and gracefully shutdown
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.Default.Fatal("failed to shutdown server", "error", err)
	}

	// shutdown queue
	if err := q.Stop(ctx); err != nil {
		logger.Default.Fatal("failed to shutdown queue", "error", err)
	}
}
