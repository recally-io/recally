package bots

import (
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/llms"
	"vibrain/internal/port/bots/handlers"

	"github.com/labstack/echo/v4"
	"gopkg.in/telebot.v3"
)

func NewReaderBot(cfg config.TelegramConfig, pool *db.Pool, e *echo.Echo, cacheService cache.Cache, llm *llms.LLM) (*Bot, error) {
	h := handlers.New(pool, llm, handlers.WithCache(cacheService))
	handlers := []Handler{
		{
			Endpoint: "/start",
			Handler: func(c telebot.Context) error {
				return c.Send(cfg.Description)
			},
			Command:     "start",
			Description: "Start the bot",
		},
		{
			Endpoint: telebot.OnText,
			Handler:  h.WebSummaryHandler,
		},
	}
	return NewBot(cfg, pool, handlers, e)
}
