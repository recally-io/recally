package handlers

import (
	"bytes"
	"io"
	"text/template"

	"recally/internal/core/bookmarks"
	"recally/internal/core/files"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader/processor"

	tele "gopkg.in/telebot.v3"
)

var imageSummaryTemplate = `
# Title
{{ .Title }}

# Description
{{ .Description }}

# Tags
{{ range $index, $tag := .Tags }}{{if $index}}, {{end}}#{{ $tag }}{{ end }}
`

var imageSummaryTempl = template.Must(template.New("imageSummaryTemplate").Parse(imageSummaryTemplate))

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

	imgUrl, err := summarier.EncodeImage(io.NopCloser(bytes.NewReader(imgBuf)), photo.FilePath)
	if err != nil {
		return c.Reply("Failed to encode image: " + err.Error())
	}

	summarier.StreamingSummary(ctx, imgUrl, streamingFunc)
	title, description, tags := summarier.ParseSummaryInfo(resp)

	var summaryBuffer bytes.Buffer
	if err := imageSummaryTempl.Execute(&summaryBuffer, struct {
		Title       string
		Description string
		Tags        []string
	}{title, description, tags}); err != nil {
		logger.FromContext(ctx).Error("failed to execute template", "err", err)

		return nil
	}

	resp = summaryBuffer.String()
	if msg, err = editMessage(c, msg, resp, true); err != nil {
		logger.FromContext(ctx).Error("failed to edit message", "err", err)

		return err
	}

	// save image content
	contentType := files.GetFileMIMEWithDefault(photo.FilePath, "image/jpeg")
	metadata := files.Metadata{
		Name:     photo.FilePath,
		Type:     photo.MediaType(),
		Ext:      files.GetFileExtensionWithDefault(photo.FilePath, "jpg"),
		Size:     photo.FileSize,
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
		Title:   title,
		Summary: description,
		S3Key:   f.S3Key,
		Tags:    tags,
	}
	if _, err = h.saveBookmark(ctx, tx, user.ID, bookmarkContent); err != nil {
		logger.FromContext(ctx).Error("failed to save bookmark", "err", err)

		return err
	}

	return nil
}
