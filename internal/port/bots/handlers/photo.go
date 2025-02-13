package handlers

import (
	"bytes"
	"fmt"
	"io"
	"recally/internal/core/bookmarks"
	"recally/internal/core/files"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader/processor"

	tele "gopkg.in/telebot.v3"
)

func (h *Handler) PhotoHandler(c tele.Context) error {
	ctx, user, tx, err := h.initHandlerRequest(c)
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
	chunkSize := 100

	streamingFunc := func(streamingMessage llms.StreamingMessage) {
		stream := streamingMessage.ToStreamingString()
		sendToUser(ctx, c, stream, &resp, &chunk, chunkSize, msg)
	}

	photo := c.Message().Photo
	if photo == nil {
		return c.Reply("Please provide a photo.")
	}

	fileReader, err := c.Bot().File(&photo.File)
	if err != nil {
		return c.Reply("Failed to download photo: " + err.Error())
	}

	imgBuf, err := io.ReadAll(fileReader)
	if err != nil {
		return c.Reply("Failed to download photo: " + err.Error())
	}

	// summarize image
	summarier := processor.NewSummaryImageProcessor(h.llm, processor.WithSummaryImageOptionUser(user))
	imgUrl, err := summarier.EncodeImage(io.NopCloser(bytes.NewReader(imgBuf)), photo.File.FilePath)
	if err != nil {
		return c.Reply("Failed to encode image: " + err.Error())
	}
	summarier.StreamingSummary(ctx, imgUrl, streamingFunc)

	// save image content
	contentType := files.GetFileMIMEWithDefault(photo.File.FilePath, "image/jpeg")
	metadata := files.Metadata{
		Name:     photo.File.FilePath,
		Type:     photo.MediaType(),
		Ext:      files.GetFileExtensionWithDefault(photo.File.FilePath, "jpg"),
		Size:     photo.File.FileSize,
		MIMEType: contentType,
	}
	f, err := files.DefaultService.UploadToS3(ctx, user.ID, "", io.NopCloser(bytes.NewReader(imgBuf)), metadata, files.WithPutObjectOptionContentType(contentType))
	if err != nil {
		return c.Reply("Failed to upload image to s3: " + err.Error())
	}

	// save bookmark
	bookmarkContent := &bookmarks.BookmarkContentDTO{
		Type:    bookmarks.ContentTypeImage,
		URL:     photo.FileURL,
		UserID:  user.ID,
		Summary: resp,
		S3Key:   f.S3Key,
	}
	bookmarkUrl, err := h.saveBookmark(ctx, tx, user.ID, bookmarkContent)
	if err == nil {
		if _, err := editMessage(c, msg, fmt.Sprintf("%s\n\n[Open Bookmark](%s)", resp, bookmarkUrl), true); err != nil {
			logger.FromContext(ctx).Error("failed to send message", "err", err)
		}
	}
	return nil
}
