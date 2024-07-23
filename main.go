package main

import (
	"context"
	"os"
	"os/signal"
	"time"
	migrations "vibrain/database"
	"vibrain/internal/core/queue"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"
	"vibrain/internal/port/bots"
	"vibrain/internal/port/httpserver"
	httpserverHandlers "vibrain/internal/port/httpserver/handlers"
)

type Service interface {
	Name() string
	Start(ctx context.Context)
	Stop(ctx context.Context)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	migrations.Migrate(ctx, config.Settings.Database.URL())

	logger.Default.Info("starting service")

	services := make([]Service, 0)

	// init basic services
	// init db pool
	pool, err := db.NewPool(ctx, config.Settings.Database.URL())
	if err != nil {
		logger.Default.Fatal("failed to create new database pool", "error", err)
	}

	// init cache service
	cacheService := cache.NewDBCache(pool)

	// start http service
	httpService, err := httpserver.New(pool, httpserverHandlers.WithCache(cacheService))
	if err != nil {
		logger.Default.Fatal("failed to create new http service", "error", err)
	}
	services = append(services, httpService)

	// start telegram bot service
	if config.Settings.Telegram.Token != "" {
		opts := make([]bots.Option, 0)
		opts = append(opts, bots.WithCache(cacheService))

		if config.Settings.Telegram.Webhook != "" {
			opts = append(opts, bots.WithWebhook(httpService.Server, config.Settings.Telegram.Webhook))
		}

		botService, err := bots.NewServer(config.Settings.Telegram.Token, pool, opts...)
		if err != nil {
			logger.Default.Fatal("failed to create new bot service", "error", err)
		}
		services = append(services, botService)
	}

	// start queue service
	queueService, err := queue.NewServer(pool)
	if err != nil {
		logger.Default.Fatal("failed to create new queue service", "error", err)
	}
	services = append(services, queueService)

	// start services
	for _, service := range services {
		go service.Start(ctx)
		logger.Default.Info("service started", "name", service.Name())
	}

	// wait for signal and gracefully shutdown
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// stop services
	for _, service := range services {
		service.Stop(ctx)
		logger.Default.Info("service stopped", "name", service.Name())
	}
}
