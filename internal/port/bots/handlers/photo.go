package handlers

import (
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader/processor"

	tele "gopkg.in/telebot.v3"
)

func (h *Handler) PhotoHandler(c tele.Context) error {
	ctx, user, _, err := h.initHandlerRequest(c)
	if err != nil {
		logger.FromContext(ctx).Error("init request error", "err", err)
		_ = c.Reply("Failed to processing message, please retry.")
		return err
	}

	msg, err := c.Bot().Reply(c.Message(), "Please wait, I'm reading the photo.")
	if err != nil {
		return err
	}

	resp := ""
	chunk := ""
	chunkSize := 400

	streamingFunc := func(streamingMessage llms.StreamingMessage) {
		stream := llms.StreamingString{
			Err: streamingMessage.Err,
		}
		if streamingMessage.Choice != nil {
			stream.Content = streamingMessage.Choice.Message.Content
		}

		msg = sendToUser(ctx, c, stream, &resp, &chunk, chunkSize, msg)
	}

	photo := c.Message().Photo
	if photo == nil {
		return c.Reply("Please provide a photo.")
	}

	summarier := processor.NewSummaryImageProcessor(h.llm, processor.WithSummaryImageOptionUser(user))
	imgUrl, err := summarier.EncodeImage(photo.FileReader, photo.File.FilePath)
	if err != nil {
		return c.Reply("Failed to encode image: " + err.Error())
	}
	summarier.StreamingSummary(ctx, imgUrl, streamingFunc)

	// bookmarkContent := &bookmarks.BookmarkContentDTO{
	// 	Type:    bookmarks.ContentTypeImage,
	// 	URL:     photo.FileURL,
	// 	UserID:  user.ID,
	// 	Summary: resp,
	// }

	// h.saveBookmark(ctx, tx, user.ID, bookmarkContent)
	return nil
}
