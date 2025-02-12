package bots

import (
	"recally/internal/core/queue"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/config"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
	"recally/internal/port/bots/handlers"

	"github.com/labstack/echo/v4"
	"gopkg.in/telebot.v3"
)

func NewReaderBot(cfg config.TelegramConfig, pool *db.Pool, e *echo.Echo, cacheService cache.Cache, llm *llms.LLM, queue *queue.Queue) (*Bot, error) {
	h := handlers.New(pool, llm, queue, handlers.WithCache(cacheService))
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
			Endpoint:    "/linkaccount",
			Handler:     h.LinkAccountHandler,
			Command:     "linkaccount",
			Description: "Link telegram bot to your account",
		},
		{
			Endpoint: telebot.OnText,
			Handler:  h.WebSummaryHandler,
		},
		{
			Endpoint: telebot.OnPhoto,
			Handler:  h.PhotoHandler,
		},
	}
	return NewBot(cfg, pool, handlers, e)
}
