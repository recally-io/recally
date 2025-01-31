package handlers

import (
	"context"
	"fmt"
	"io"
	"recally/internal/core/bookmarks"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader/fetcher"
	"recally/internal/pkg/webreader/processor"
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
	logger.FromContext(ctx).Info("TextHandler start summary", "text", text)
	url := getUrlFromText(text)
	if url == "" {
		return c.Reply("Please provide a valid URL.")
	}

	processSendError := func(err error) error {
		logger.FromContext(ctx).Error("TextHandler failed to send message", "err", err, "text", text)
		err = c.Reply("Failed to send message. " + err.Error())
		if err != nil {
			logger.FromContext(ctx).Error("error reply message", "err", err)
		}
		return err
	}

	msg, err := c.Bot().Reply(c.Message(), "Please wait, I'm reading the page.")
	if err != nil {
		return processSendError(err)
	}

	resp := ""
	chunk := ""
	chunkSize := 400
	isSummaryCached := true

	editMessage := func(msg *tele.Message, text string, format bool) (*tele.Message, error) {
		if format {
			return c.Bot().Edit(msg, convertToTGMarkdown(text), tele.ModeMarkdownV2)
		}
		return c.Bot().Edit(msg, text)
	}

	sendToUser := func(stream llms.StreamingString) {
		line, err := stream.Content, stream.Err
		chunk += line

		if err != nil {
			if err == io.EOF {
				resp += chunk
				resp = strings.ReplaceAll(resp, "\\n", "\n")
				if _, err := editMessage(msg, resp, true); err != nil {
					if strings.Contains(err.Error(), "message is not modified") {
						return
					}
					_ = processSendError(err)
				}
				return
			}
			logger.FromContext(ctx).Error("TextHandler failed to get summary", "err", err)
			if msg, err = editMessage(msg, "Failed to get summary.", false); err != nil {
				_ = processSendError(err)
				return
			}
		}

		if len(chunk) > chunkSize {
			resp += chunk
			chunk = ""
			var newErr error
			msg, newErr = editMessage(msg, resp, false)
			if newErr != nil {
				_ = processSendError(err)
				return
			}
		}
	}

	var bookmarkContentDTO *bookmarks.BookmarkContentDTO
	// cache the summary
	summary, err := cache.RunInCache[string](ctx, cache.DefaultDBCache, cache.NewCacheKey("WebSummary", url), 24*time.Hour, func() (*string, error) {
		isSummaryCached = false
		// cache the content
		content, err := h.bookmarkService.FetchWebContentWithCache(ctx, url, fetcher.FetchOptions{
			FecherType: fetcher.TypeJinaReader,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get content: %w", err)
		}
		bookmarkContentDTO.FromReaderContent(content)

		// process the summary
		summarier := processor.NewSummaryProcessor(h.llm, processor.WithSummaryOptionUser(user))
		summarier.StreamingSummary(ctx, content.Markwdown, sendToUser)
		return &resp, nil
	})
	if err != nil {
		return processSendError(err)
	}

	// if summary is cached, just return the cached summary
	// if not cached, streaming send summary to user and save the bookmark
	if isSummaryCached {
		if _, err := editMessage(msg, *summary, true); err != nil {
			logger.FromContext(ctx).Error("TextHandler failed to send message", "err", err, "text", text)
		}
	} else {
		bookmarkContentDTO.Summary = resp
		h.saveBookmark(ctx, tx, user.ID, bookmarkContentDTO)
	}
	return nil
}

func getUrlFromText(text string) string {
	return urlPattern.FindString(text)
}

func (h *Handler) saveBookmark(ctx context.Context, tx db.DBTX, userId uuid.UUID, bookmarkContent *bookmarks.BookmarkContentDTO) {
	bookmark, err := h.bookmarkService.CreateBookmark(ctx, tx, userId, bookmarkContent)
	if err != nil {
		logger.FromContext(ctx).Error("save bookmark from reader bot error", "err", err.Error())
	} else {
		logger.FromContext(ctx).Info("save bookmark from reader bot", "id", bookmark.ID)
	}
}
