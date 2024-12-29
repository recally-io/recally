package main

import (
	"context"
	"os"
	"os/signal"
	migrations "recally/database"
	"recally/internal/core/queue"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/config"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/s3"
	"recally/internal/port/bots"
	"recally/internal/port/httpserver"
	"time"
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

	llm := llms.New(config.Settings.OpenAI.BaseURL, config.Settings.OpenAI.ApiKey)
	s3Client, err := s3.New(config.Settings.S3)
	if err != nil {
		logger.Default.Fatal("failed to create new s3 client", "error", err)
	}

	// start queue service
	queueService, err := queue.NewServer(pool, llm)
	if err != nil {
		logger.Default.Fatal("failed to create new queue service", "error", err)
	}
	services = append(services, queueService)

	// start http service
	httpService, err := httpserver.New(pool, llm, queueService.Queue, httpserver.WithCache(cacheService), httpserver.WithS3(s3Client))
	if err != nil {
		logger.Default.Fatal("failed to create new http service", "error", err)
	}
	services = append(services, httpService)

	// start telegram bot service
	if config.Settings.Telegram.Reader.Token != "" {
		cfg := config.Settings.Telegram.Reader
		botService, err := bots.NewServer(bots.ReaderBot, cfg, pool, httpService.Server, cacheService, llm, queueService.Queue)
		if err != nil {
			logger.Default.Fatal("failed to create new bot service", "error", err, "type", bots.ReaderBot, "name", cfg.Name)
		}
		services = append(services, botService)
	}

	if config.Settings.Telegram.Chat.Token != "" {
		cfg := config.Settings.Telegram.Chat
		botService, err := bots.NewServer(bots.ChatBot, cfg, pool, httpService.Server, cacheService, llm, queueService.Queue)
		if err != nil {
			logger.Default.Fatal("failed to create new bot service", "error", err, "type", bots.ChatBot, "name", cfg.Name)
		}
		services = append(services, botService)
	}

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
