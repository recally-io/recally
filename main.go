package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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

	migrations "recally/database"
)

// Build information injected via ldflags.
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
	builtBy = "unknown"
)

type Service interface {
	Name() string
	Start(ctx context.Context)
	Stop(ctx context.Context)
}

func main() {
	// Handle version and health commands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Printf("Recally %s\n", version)
			fmt.Printf("Commit: %s\n", commit)
			fmt.Printf("Date: %s\n", date)
			fmt.Printf("Built by: %s\n", builtBy)

			return
		case "health":
			// Simple health check - could be enhanced to check database connectivity
			fmt.Println("OK")

			return
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	migrations.Migrate(ctx, config.Settings.Database.URL())

	logger.Default.Info("starting service", "version", version, "commit", commit)

	services := make([]Service, 0)

	// init basic services using default config
	pool := db.DefaultPool
	cacheService := cache.DefaultDBCache
	llm := llms.DefaultLLM
	s3Client := s3.DefaultClient
	riverQueue := queue.DefaultQueue

	// start queue service
	queueService, err := queue.NewServer(riverQueue)
	if err != nil {
		logger.Default.Fatal("failed to create new queue service", "err", err)
	}

	services = append(services, queueService)

	// start http service
	httpService, err := httpserver.New(pool, llm, queueService.Queue, httpserver.WithCache(cacheService), httpserver.WithS3(s3Client))
	if err != nil {
		logger.Default.Fatal("failed to create new http service", "err", err)
	}

	services = append(services, httpService)

	// start telegram bot service
	if config.Settings.Telegram.Reader.Token != "" {
		cfg := config.Settings.Telegram.Reader
		botService, err := bots.NewServer(bots.ReaderBot, cfg, pool, httpService.Server, cacheService, llm, queueService.Queue)
		if err != nil {
			logger.Default.Fatal("failed to create new bot service", "err", err, "type", bots.ReaderBot, "name", cfg.Name)
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
