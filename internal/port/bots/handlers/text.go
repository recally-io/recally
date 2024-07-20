package handlers

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
	"vibrain/internal/core/workers"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/constant"
	"vibrain/internal/pkg/logger"

	tele "gopkg.in/telebot.v3"
)

var urlPattern = regexp.MustCompile(`http[s]?://(?:[-\w.]|(?:%[\da-fA-F]{2}))+`)

func (h *Handler) TextHandler(c tele.Context) error {
	user := c.Sender()
	text := c.Text()
	ctx := c.Get(constant.ContextKeyContext).(context.Context)
	logger.FromContext(ctx).Info("TextHandler", "user", user.Username, "text", text)
	url := getUrlFromText(text)
	if url == "" {
		return c.Reply("Please provide a valid URL.")
	}
	reader, err := h.worker.WebSummaryStream(ctx, url)
	if err != nil {
		return c.Reply(fmt.Sprintf("Failed to get summary:\n%s", err.Error()))
	}
	defer reader.Close()

	processSendError := func(err error) error {
		logger.FromContext(ctx).Error("TextHandler failed to send message", "error", err.Error())
		return c.Reply("Failed to send message. " + err.Error())
	}

	msg, err := c.Bot().Send(user, "Please wait, I'm reading the page.")
	if err != nil {
		return processSendError(err)
	}
	resp := ""
	chunk := ""
	chunkSize := 400
	for {
		line, err := reader.Stream()
		chunk += line

		if err != nil {
			if err == io.EOF {
				resp += chunk
				resp = strings.ReplaceAll(resp, "\\n", "\n")
				if h.Cache != nil {
					cacheKey := cache.NewCacheKey(workers.WebSummaryCacheDomian, url)
					h.Cache.SetWithContext(ctx, cacheKey, resp, 24*time.Hour)
				}
				if _, err := c.Bot().Edit(msg, convertToTGMarkdown(resp), tele.ModeMarkdownV2); err != nil {
					if strings.Contains(err.Error(), "message is not modified") {
						return nil
					}
					return processSendError(err)
				}
				return nil
			}
			logger.FromContext(ctx).Error("TextHandler", "error", err.Error())
			if _, err := c.Bot().Edit(msg, "Failed to get summary."); err != nil {
				return processSendError(err)
			}
		}

		if len(chunk) > chunkSize {
			resp += chunk
			chunk = ""
			var newErr error
			msg, newErr = c.Bot().Edit(msg, resp)
			if newErr != nil {
				return processSendError(err)
			}
		}
	}
}

func getUrlFromText(text string) string {
	return urlPattern.FindString(text)
}
