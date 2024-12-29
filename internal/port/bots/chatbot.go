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

func NewChatBot(cfg config.TelegramConfig, pool *db.Pool, e *echo.Echo, cacheService cache.Cache, llm *llms.LLM, queue *queue.Queue) (*Bot, error) {
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
			Endpoint:    "/new_assistant",
			Handler:     h.LLMChatNewAssistanthandler,
			Command:     "new_assistant",
			Description: "create a new assistant",
		},
		{
			Endpoint:    "/list_assistants",
			Handler:     h.LLMChatListAssistantshandler,
			Command:     "list_assistants",
			Description: "list all assistants",
		},
		{
			Endpoint:    "/new_thread",
			Handler:     h.LLMChatNewThreadHandler,
			Command:     "new_thread",
			Description: "create a new conversation",
		},
		{
			Endpoint:    "/list_threads",
			Handler:     h.LLMChatListThreadHandler,
			Command:     "list_threads",
			Description: "list all conversation threads for the assistant",
		},
		{
			Endpoint: telebot.OnText,
			Handler:  h.LLMChatHandler,
		},
	}
	return NewBot(cfg, pool, handlers, e)
}
