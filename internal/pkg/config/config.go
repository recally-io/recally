package config

import (
	"vibrain/internal/pkg/logger"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var Settings = &Config{}

// type DatabaseConfig struct {
// 	Driver   string `json:"driver"` // postgres, mysql, sqlite3
// 	Host     string `json:"host"`
// 	Port     int    `json:"port"`
// 	User     string `json:"user"`
// 	Password string `json:"password"`
// 	Name     string `json:"name"`
// }

// func (d DatabaseConfig) DSN() string {
// 	switch d.Driver {
// 	case "postgres":
// 		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", d.Host, d.Port, d.User, d.Password, d.Name)
// 	case "mysql":
// 		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.User, d.Password, d.Host, d.Port, d.Name)
// 	case "sqlite3":
// 		return d.Name
// 	}
// 	return ""
// }

type Config struct {
	Debug            bool   `env:"DEBUG" envDefault:"false"`
	Port             int    `env:"PORT" envDefault:"1323"`

	DatabaseURL      string `env:"DATABASE_URL,required"`
	QueueDatabaseURL string `env:"QUEUE_DATABASE_URL,expand" envDefault:"${DATABASE_URL}"`

	TelegramToken string `env:"TELEGRAM_TOKEN,required"`
}

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Default.Fatal("Error loading .env file", "err", err)
	}
	logger.Default.Info(".env file loaded")

	// Load env to settings
	if err := env.Parse(Settings); err != nil {
		logger.Default.Fatal("failed to load settings", "err", err)
	}
	logger.Default.Info("settings loaded")
}
