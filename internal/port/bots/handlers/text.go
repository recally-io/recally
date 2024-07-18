package handlers

import (
	"context"
	"regexp"
	"vibrain/internal/core/workers"
	"vibrain/internal/pkg/constant"
	"vibrain/internal/pkg/logger"

	tele "gopkg.in/telebot.v3"
)

var urlPattern = regexp.MustCompile(`http[s]?://(?:[-\w.]|(?:%[\da-fA-F]{2}))+`)

func TextHandler(c tele.Context) error {
	user := c.Sender()
	text := c.Text()
	ctx := c.Get(constant.ContextKeyContext).(context.Context)
	logger.FromContext(ctx).Info("TextHandler", "user", user.Username, "text", text)
	url := getUrlFromText(text)
	if url == "" {
		return c.Reply("Please provide a valid URL.")
	}
	summary, err := workers.WebSummary(ctx, url)
	if err != nil {
		return c.Reply("Failed to get summary.")
	}
	return c.Reply(summary)
}

func getUrlFromText(text string) string {
	return urlPattern.FindString(text)
}
