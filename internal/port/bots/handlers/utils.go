package handlers

import (
	"context"
	"io"
	"regexp"
	"strings"

	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"

	tele "gopkg.in/telebot.v3"
)

var urlPattern = regexp.MustCompile(`http[s]?://(?:[-\w.]|(?:%[\da-fA-F]{2}))+/?([^\s]*)`)

func getUrlFromText(text string) string {
	return urlPattern.FindString(text)
}

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

func sendToUser(ctx context.Context, c tele.Context, stream llms.StreamingString, resp, chunk *string, chunkSize int, msg *tele.Message) {
	line, err := stream.Content, stream.Err
	if line == "" && err == nil {
		return
	}

	*chunk += line

	if err != nil {
		if err == io.EOF {
			*resp += *chunk
			*resp = strings.ReplaceAll(*resp, "\\n", "\n")
			// if _, err := editMessage(c, msg, *resp, true); err != nil {
			// 	if strings.Contains(err.Error(), "message is not modified") {
			// 		return
			// 	}
			// 	_ = processSendError(ctx, c, err)
			// }
			return
		}

		logger.FromContext(ctx).Error("TextHandler failed to get summary", "err", err)

		if _, err = editMessage(c, msg, "Failed to get summary.", false); err != nil {
			_ = processSendError(ctx, c, err)

			return
		}
	}

	if len(*chunk) > chunkSize {
		*resp += *chunk
		*chunk = ""

		if newMsg, newErr := editMessage(c, msg, *resp, false); newErr != nil {
			_ = processSendError(ctx, c, newErr)

			return
		} else {
			*msg = *newMsg
		}
	}
}
