package config

import (
	"os"
	"vibrain/internal/pkg/logger"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var Settings = &Config{}

type OAuthConfig struct {
	Provider string   `env:"PROVIDER,required"`
	Endpoint string   `env:"ENDPOINT"`
	Key      string   `env:"KEY,required"`
	Secret   string   `env:"SECRET,required"`
	Scopes   []string `env:"SCOPES"`
}

type ServiceConfig struct {
	Fqdn string `env:"FQDN" envDefault:"localhost:1323"`
	Port int    `env:"PORT" envDefault:"1323"`
}

type Config struct {
	Debug   bool          `env:"DEBUG" envDefault:"false"`
	Service ServiceConfig `envPrefix:"SERVICE"`

	DatabaseURL      string `env:"DATABASE_URL,required"`
	QueueDatabaseURL string `env:"QUEUE_DATABASE_URL,expand" envDefault:"${DATABASE_URL}"`

	TelegramToken string `env:"TELEGRAM_TOKEN"`

	JWTSecret string        `env:"JWT_SECRET,required"`
	OAuths    []OAuthConfig `envPrefix:"OAUTH"`
}

func init() {
	// Load .env file if exists
	if _, err := os.Stat(".env"); err == nil {
		logger.Default.Info("loading .env file")
		if err := godotenv.Load(); err != nil {
			logger.Default.Fatal("Error loading .env file", "err", err)
		}
		logger.Default.Info(".env file loaded")
	} else {
		logger.Default.Info(".env file not found")
	}
	// Load env to settings
	if err := env.Parse(Settings); err != nil {
		logger.Default.Fatal("failed to load settings", "err", err)
	}
	logger.Default.Info("settings loaded")
}
