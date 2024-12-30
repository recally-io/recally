package handlers

import (
	"context"
	"fmt"
	"io"
	"recally/internal/core/bookmarks"
	"recally/internal/core/queue"
	"recally/internal/core/workers"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	tele "gopkg.in/telebot.v3"
)

var urlPattern = regexp.MustCompile(`http[s]?://(?:[-\w.]|(?:%[\da-fA-F]{2}))+/?([^\s]*)`)

func (h *Handler) WebSummaryHandler(c tele.Context) error {
	ctx, user, tx, err := h.initHandlerRequest(c)
	if err != nil {
		logger.FromContext(ctx).Error("init request error", "err", err)
		_ = c.Reply("Failed to processing message, please retry.")
		return err
	}

	text := c.Text()
	logger.FromContext(ctx).Info("TextHandler", "user", user.Username, "text", text)
	url := getUrlFromText(text)
	if url == "" {
		return c.Reply("Please provide a valid URL.")
	}
	reader, err := h.toolService.WebSummaryStream(ctx, url)
	if err != nil {
		return c.Reply(fmt.Sprintf("Failed to get summary:\n%s", err.Error()))
	}
	defer reader.Close()

	processSendError := func(err error) error {
		logger.FromContext(ctx).Error("TextHandler failed to send message", "error", err.Error())
		return c.Reply("Failed to send message. " + err.Error())
	}

	msg, err := c.Bot().Send(c.Sender(), "Please wait, I'm reading the page.")
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
				cacheKey := cache.NewCacheKey(workers.WebSummaryCacheDomian, url)
				h.cache.SetWithContext(ctx, cacheKey, resp, 24*time.Hour)
				if _, err := c.Bot().Edit(msg, convertToTGMarkdown(resp), tele.ModeMarkdownV2); err != nil {
					if strings.Contains(err.Error(), "message is not modified") {
						return nil
					}
					return processSendError(err)
				}
				h.saveBookmark(ctx, tx, url, user.ID, resp)
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

func (h *Handler) saveBookmark(ctx context.Context, tx db.DBTX, url string, userId uuid.UUID, summary string) {
	bookmark, err := h.bookmarkService.Create(ctx, tx, &bookmarks.BookmarkDTO{
		UserID:    userId,
		URL:       url,
		Summary:   summary,
		CreatedAt: time.Now(),
	})
	if err != nil {
		logger.FromContext(ctx).Error("save bookmark from reader bot error", "err", err.Error())
	} else {
		logger.FromContext(ctx).Info("save bookmark from reader bot", "id", bookmark.ID, "title", bookmark.Title)
	}

	result, err := h.queue.Insert(ctx, queue.CrawlerWorkerArgs{
		ID:          bookmark.ID,
		UserID:      bookmark.UserID,
		FetcherName: bookmarks.JinaFetcher,
	}, nil)
	if err != nil {
		logger.FromContext(ctx).Error("failed to insert job", "err", err)
	} else {
		logger.FromContext(ctx).Info("success inserted job", "result", result, "err", err)
	}
}
