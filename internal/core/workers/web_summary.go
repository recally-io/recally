package workers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"vibrain/internal/pkg/logger"
)

const elmoSummaryUrl = "https://www.elmo.chat/api/v1/prompt"

type StreamStringReader struct {
	io.Closer
	*bufio.Reader
	Text string
}

func (r *StreamStringReader) Stream() (string, error) {
	if r.Text != "" {
		return r.Text, io.EOF
	}
	for {
		// Read until the next newline
		line, err := r.ReadBytes('\n')
		if err != nil {
			// If the error is EOF, send it to errChan and return
			if errors.Is(err, io.EOF) {
				return "", io.EOF
			}
			return "", err
		}
		if bytes.HasPrefix(line, []byte("0:")) {
			// using regex to extract the summary text, origin text will be '0:"som text"', and I only want 'some text'
			text := string(line[3 : len(line)-2])
			return strings.ReplaceAll(text, "\\n", "\n"), nil
		}
	}
}

func (r *StreamStringReader) Close() {
	if r.Closer != nil {
		r.Closer.Close()
	}
}

func (w *Worker) WebSummary(ctx context.Context, url string) (string, error) {
	cacheKey := fmt.Sprintf("WebSummary:%s", url)
	reader, err := w.WebSummaryStream(ctx, url)
	if err != nil {
		return "", err
	}
	sb := strings.Builder{}
	for {
		line, err := reader.Stream()
		if err != nil {
			if errors.Is(err, io.EOF) {
				summary := sb.String()
				summary = strings.ReplaceAll(summary, "\\n", "\n")
				if w.cache != nil {
					w.cache.SetWithContext(ctx, cacheKey, summary, 24*time.Hour)
				}
				return summary, nil
			}
			return "", err
		}
		sb.WriteString(line)
	}
}

func (w *Worker) WebSummaryStream(ctx context.Context, url string) (*StreamStringReader, error) {
	return w.elmoSummary(ctx, url, "")
}

func (w *Worker) elmoSummary(ctx context.Context, url, pageContent string) (*StreamStringReader, error) {
	// get result from cache
	cacheKey := fmt.Sprintf("WebSummary:%s", url)
	reader := &StreamStringReader{}
	if w.cache != nil {
		if val, ok := w.cache.GetWithContext(ctx, cacheKey); ok {
			logger.FromContext(ctx).Info("WebSummary", "cache", "hit", "url", url)
			reader = &StreamStringReader{
				Text: val.(string),
			}
			return reader, nil
		}
	}
	if pageContent == "" {
		content, err := w.WebReader(ctx, url)
		if err != nil {
			return reader, fmt.Errorf("failed to get content: %w", err)
		}
		pageContent = content.Content
		url = content.Url
	}

	body, _ := json.Marshal(map[string]interface{}{
		"regenerate":  false,
		"enableCache": true,
		"conversation": []map[string]string{
			{
				"role":    "user",
				"content": "/summarize",
			},
		},
		"metadata": map[string]map[string]string{
			"system":  {"language": "zh-Hans"},
			"website": {"url": url, "content": pageContent},
		},
	})

	req, err := http.NewRequest("POST", elmoSummaryUrl, bytes.NewBuffer(body))
	if err != nil {
		return reader, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

	client := newHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return reader, fmt.Errorf("failed to send request: %w", err)
	}
	// defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respData, err := io.ReadAll(resp.Body)
		if err != nil {
			return reader, fmt.Errorf("failed to read response body: %w", err)
		}
		slog.Error("summary url error", "url", url, "status", resp.Status, "err", string(respData))
		return reader, fmt.Errorf("failed to get summary: %s, %v", resp.Status, string(respData))
	}
	reader = &StreamStringReader{
		Closer: resp.Body,
		Reader: bufio.NewReader(resp.Body),
	}
	return reader, nil
}
