package config

import (
	"fmt"
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
	Host string `env:"HOST" envDefault:"localhost"`
	Port int    `env:"PORT" envDefault:"1323"`
}

type TelegramConfig struct {
	Token       string `env:"TOKEN,required"`
	Name        string `env:"NAME"`
	Description string `env:"DESCRIPTION"`
	Webhook     bool   `env:"WEBHOOK" envDefault:"false"`
}

type DatabaseConfig struct {
	Driver   string `env:"DRIVER" envDefault:"postgres"`
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     int    `env:"PORT" envDefault:"5432"`
	User     string `env:"USER" envDefault:"postgres"`
	Password string `env:"PASSWORD" envDefault:"postgres"`
	Database string `env:"DATABASE" envDefault:"postgres"`
}

func (db DatabaseConfig) URL() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable", db.Driver, db.User, db.Password, db.Host, db.Port, db.Database)
}

type Config struct {
	Debug   bool          `env:"DEBUG" envDefault:"false"`
	Service ServiceConfig `envPrefix:"SERVICE_"`

	Database DatabaseConfig `envPrefix:"DATABASE_"`
	Telegram TelegramConfig `envPrefix:"TELEGRAM_"`

	JWTSecret string        `env:"JWT_SECRET,required"`
	OAuths    []OAuthConfig `envPrefix:"OAUTH_"`
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
