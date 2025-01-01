package config

import (
	"fmt"
	"os"
	"recally/internal/pkg/logger"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var Settings = &Config{}

type OAuthConfig struct {
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
	Name                string `env:"NAME" envDefault:"RecallyBot"`
	Token               string `env:"TOKEN"`
	Webhook             bool   `env:"WEBHOOK"`
	WebhookSecrectToken string `env:"WEBHOOK_SECRET_TOKEN"`
	Description         string `env:"DESCRIPTION" envDefault:"Hi, I'am Recally Bot. Contact me for more information on Twitter https://x.com/LiuVaayne"`
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

type OpenAIConfig struct {
	BaseURL string `env:"BASE_URL" envDefault:"https://api.openai.com"`
	ApiKey  string `env:"API_KEY,required"`
	Model   string `env:"MODEL" envDefault:"gpt-4o-mini"`
}

type GoogleSearchConfig struct {
	ApiKey   string `env:"API_KEY"`
	EngineID string `env:"ENGINE_ID"`
}

type S3Config struct {
	Endpoint        string `env:"ENDPOINT"`
	AccessKeyID     string `env:"ACCESS_KEY_ID"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY"`
	BucketName      string `env:"BUCKET_NAME"`
	PublicURL       string `env:"PUBLIC_URL"`
}

type Config struct {
	// Env is the environment the service is running in
	// such as dev, staging, production
	Env     string        `env:"ENV" envDefault:"production"`
	Debug   bool          `env:"DEBUG" envDefault:"false"`
	DebugUI bool          `env:"DEBUG_UI" envDefault:"false"`
	Service ServiceConfig `envPrefix:"SERVICE_"`

	Database DatabaseConfig `envPrefix:"DATABASE_"`
	Telegram struct {
		Reader  TelegramConfig `envPrefix:"READER_"`
		Chat    TelegramConfig `envPrefix:"CHAT_"`
		MemChat TelegramConfig `envPrefix:"MEMCHAT_"`
	} `envPrefix:"TELEGRAM_"`

	JWTSecret string `env:"JWT_SECRET,required"`
	OAuths    struct {
		Github OAuthConfig `envPrefix:"GITHUB_"`
	} `envPrefix:"OAUTH_"`
	OpenAI            OpenAIConfig       `envPrefix:"OPENAI_"`
	GoogleSearch      GoogleSearchConfig `envPrefix:"GOOGLE_SEARCH_"`
	S3                S3Config           `envPrefix:"S3_"`
	BrowserControlUrl string             `env:"BROWSER_CONTROL_URL" envDefault:"http://localhost:9222"`
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
