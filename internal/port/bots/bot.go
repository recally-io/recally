package bots

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Bot struct {
	*telebot.Bot
	cfg         config.TelegramConfig
	handlers    []Handler
	webhookPath string
}

type BotType string

const (
	ReaderBot BotType = "readerbot"
	ChatBot   BotType = "chatbot"
)

type BotOption func(Bot)

type Handler struct {
	// Endpoint is the endpoint that will be used to handle the command
	Endpoint any
	// Handler is the function that will be called when the command is received
	Handler func(c telebot.Context) error
	// Command is the command that will be shown in the list of commands
	Command string
	// Description is the description of the command
	Description string
}

func NewBot(cfg config.TelegramConfig, pool *db.Pool, handlers []Handler, e *echo.Echo) (*Bot, error) {
	bot := &Bot{
		cfg: cfg,
	}

	b, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new bot: %w", err)
	}
	registerMiddlewarw(b, pool)

	for _, handler := range handlers {
		b.Handle(handler.Endpoint, handler.Handler)
	}
	bot.handlers = handlers
	bot.Bot = b

	if e != nil && cfg.Webhook {
		bot.webhookPath = fmt.Sprintf("/telegram/bot/%s/%s", cfg.Name, uuid.New().String())
		setWebhook(bot, e, bot.webhookPath)
	}
	return bot, nil
}

func registerMiddlewarw(b *telebot.Bot, db *db.Pool) {
	b.Use(contextMiddleware())
	b.Use(middleware.Recover())
	b.Use(TransactionMiddleware(db))
}

func setWebhook(b *Bot, e *echo.Echo, webhookPath string) {
	e.POST(webhookPath, func(c echo.Context) error {
		if b.cfg.WebhookSecrectToken != "" && c.Request().Header.Get("X-Telegram-Bot-Api-Secret-Token") != b.cfg.WebhookSecrectToken {
			logger.FromContext(c.Request().Context()).Error("invalid secret token in request")
			return c.String(http.StatusUnauthorized, "invalid secret token")
		}

		var update telebot.Update
		if err := json.NewDecoder(c.Request().Body).Decode(&update); err != nil {
			logger.FromContext(c.Request().Context()).Error("cannot decode update", "err", err)
			return c.String(http.StatusBadRequest, fmt.Sprintf("cannot decode update: %s", err))
		}
		b.Updates <- update
		return nil
	})
}

func (b *Bot) Start(ctx context.Context) {
	logger := logger.Default
	if err := b.SetMyName(b.cfg.Name, "en"); err != nil {
		logger.Error("failed to set bot name", "err", err)
	} else {
		logger.Info("success set bot name", "name", b.cfg.Name)
	}
	if err := b.SetMyDescription(b.cfg.Description, "en"); err != nil {
		logger.Error("failed to set bot description", "err", err)
	} else {
		logger.Info("success set bot description", "description", b.cfg.Description)
	}
	commands := b.getCommands()
	if len(commands) > 0 {
		if err := b.SetCommands(commands); err != nil {
			logger.Error("failed to set commands", "err", err)
		} else {
			logger.Info("success set commands", "commands", commands)
		}
	}

	if b.webhookPath != "" {
		params := map[string]string{
			"url":                  fmt.Sprintf("%s%s", config.Settings.Service.Fqdn, b.webhookPath),
			"drop_pending_updates": "true",
		}
		if b.cfg.WebhookSecrectToken != "" {
			params["secret_token"] = b.cfg.WebhookSecrectToken
		}
		if _, err := b.Raw("setWebhook", params); err != nil {
			logger.Error("failed to set webhook", "err", err)
		} else {
			logger.Info("success set webhook", "url", params["url"])
		}
		return
	}
	b.Bot.Start()
}

func (b *Bot) Stop(ctx context.Context) {
	if b.webhookPath != "" {
		dropPending := true
		// RemoveWebhook
		if _, err := b.Raw("deleteWebhook", map[string]bool{
			"drop_pending_updates": dropPending,
		}); err != nil {
			logger.FromContext(ctx).Error("failed to remove webhook", "err", err)
		} else {
			logger.FromContext(ctx).Info("success remove webhook", "drop_pending_updates", dropPending)
		}
		return
	}
	b.Bot.Stop()
}

func (b *Bot) AddHandler(handler Handler) {
	b.handlers = append(b.handlers, handler)
}

func (b *Bot) AddHandlers(handlers []Handler) {
	b.handlers = append(b.handlers, handlers...)
}

func (b *Bot) getCommands() []telebot.Command {
	var commands []telebot.Command
	for _, handler := range b.handlers {
		if handler.Command != "" {
			commands = append(commands, telebot.Command{
				Text:        handler.Command,
				Description: handler.Description,
			})
		}
	}
	return commands
}
