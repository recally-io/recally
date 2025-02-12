package handlers

import (
	"context"
	"io"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func processSendError(ctx context.Context, c tele.Context, err error) error {
	logger.FromContext(ctx).Error("failed to send message", "err", err)
	err = c.Reply("Failed to send message. " + err.Error())
	if err != nil {
		logger.FromContext(ctx).Error("error reply message", "err", err)
	}
	return err
}

func editMessage(c tele.Context, msg *tele.Message, text string, format bool) (*tele.Message, error) {
	if format {
		return c.Bot().Edit(msg, convertToTGMarkdown(text), tele.ModeMarkdownV2)
	}
	return c.Bot().Edit(msg, text)
}

func sendToUser(ctx context.Context, c tele.Context, stream llms.StreamingString, resp, chunk *string, chunkSize int, msg *tele.Message) *tele.Message {
	line, err := stream.Content, stream.Err
	*chunk += line

	if err != nil {
		if err == io.EOF {
			*resp += *chunk
			*resp = strings.ReplaceAll(*resp, "\\n", "\n")
			if _, err := editMessage(c, msg, *resp, true); err != nil {
				if strings.Contains(err.Error(), "message is not modified") {
					return msg
				}
				_ = processSendError(ctx, c, err)
			}
			return msg
		}
		logger.FromContext(ctx).Error("TextHandler failed to get summary", "err", err)
		if _, err = editMessage(c, msg, "Failed to get summary.", false); err != nil {
			_ = processSendError(ctx, c, err)
			return nil
		}
	}

	if len(*chunk) > chunkSize {
		*resp += *chunk
		*chunk = ""
		var newErr error
		msg, newErr = editMessage(c, msg, *resp, false)
		if newErr != nil {
			_ = processSendError(ctx, c, newErr)
			return msg
		}
	}
	return msg
}
