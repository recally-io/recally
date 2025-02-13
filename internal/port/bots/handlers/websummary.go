package handlers

import (
	"context"
	"fmt"
	"recally/internal/core/bookmarks"
	"recally/internal/core/queue"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/config"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader/fetcher"
	"recally/internal/pkg/webreader/processor"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
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

	msg, err := c.Bot().Reply(c.Message(), "Please wait, I'm reading the page.")
	if err != nil {
		return processSendError(ctx, c, err)
	}

	resp := ""
	chunk := ""
	chunkSize := 400

	sendToUser := func(stream llms.StreamingString) {
		sendToUser(ctx, c, stream, &resp, &chunk, chunkSize, msg)
	}

	var bookmarkContentDTO bookmarks.BookmarkContentDTO
	// cache the summary
	summary, err := cache.RunInCache[string](ctx, cache.DefaultDBCache, cache.NewCacheKey("WebSummary", url), 24*time.Hour, func() (*string, error) {
		// cache the content
		content, err := h.bookmarkService.FetchWebContentWithCache(ctx, url, fetcher.FetchOptions{
			FecherType: fetcher.TypeHttp,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get content: %w", err)
		}
		bookmarkContentDTO.FromReaderContent(content)
		bookmarkContentDTO.Type = bookmarks.ContentTypeBookmark

		// process the summary
		summarier := processor.NewSummaryProcessor(h.llm, processor.WithSummaryOptionUser(user))
		summarier.StreamingSummary(ctx, content.Markwdown, sendToUser)
		return &resp, nil
	})
	if err != nil {
		return processSendError(ctx, c, err)
	}
	var tags []string
	*summary, tags = processor.NewSummaryProcessor(h.llm).ParseSummaryInfo(*summary)
	bookmarkContentDTO.Summary = *summary
	bookmarkContentDTO.Tags = append(bookmarkContentDTO.Tags, tags...)
	if msg, err = editMessage(c, msg, *summary, true); err != nil {
		return err
	}

	if _, err = h.saveBookmark(ctx, tx, user.ID, &bookmarkContentDTO); err != nil {
		logger.FromContext(ctx).Error("failed to save bookmark", "err", err)
		return err
	}
	return nil
}

func getUrlFromText(text string) string {
	return urlPattern.FindString(text)
}

func (h *Handler) saveBookmark(ctx context.Context, tx pgx.Tx, userId uuid.UUID, bookmarkContent *bookmarks.BookmarkContentDTO) (string, error) {
	bookmark, err := h.bookmarkService.CreateBookmark(ctx, tx, userId, bookmarkContent)
	if err != nil {
		logger.FromContext(ctx).Error("save bookmark from reader bot error", "err", err.Error())
		return "", err
	} else {
		logger.FromContext(ctx).Info("save bookmark from reader bot", "id", bookmark.ID)
	}

	result, err := queue.DefaultQueue.InsertTx(ctx, tx, queue.CrawlerWorkerArgs{
		ID:           bookmark.ID,
		UserID:       bookmark.UserID,
		FetchOptions: fetcher.FetchOptions{FecherType: fetcher.TypeHttp},
	}, &river.InsertOpts{
		ScheduledAt: time.Now().Add(time.Second * 5),
	})
	if err != nil {
		logger.FromContext(ctx).Error("failed to insert job", "err", err)
	} else {
		logger.FromContext(ctx).Info("success inserted job", "result", result, "err", err)
	}
	bookmarkUrl := fmt.Sprintf("%s/bookmarks/%s", config.Settings.Service.Fqdn, bookmark.ID)
	return bookmarkUrl, nil
}
