package jinareader

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"vibrain/internal/pkg/logger"
	"vibrain/internal/pkg/tools"
)

const (
	jinaHost = "https://r.jina.ai"
)

type Tool struct {
	tools.BaseTool
	httpClient *http.Client
}

type RequestArgs struct {
	Url string `json:"url" jsonschema_description:"The URL to read."`
}

type Content struct {
	Url         string `json:"url"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
}

type Result struct {
	Code   int     `json:"code"`
	Status float64 `json:"status"`
	Data   Content `json:"data"`
}

func New() *Tool {
	return &Tool{
		BaseTool: tools.BaseTool{
			Name:        "jinareader",
			Description: "Get LLM-friendly content from a URL.",
			Parameters:  &RequestArgs{},
		},
		httpClient: &http.Client{
			Timeout: 5 * 60 * time.Second,
		},
	}
}

func (t *Tool) Invoke(ctx context.Context, args string) (string, error) {
	var params RequestArgs
	if err := t.UnmarshalArgs(ctx, args, &params); err != nil {
		return "", err
	}

	result, err := t.Read(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to invoke tool: %w", err)
	}
	return t.MarshalResult(ctx, result)
}

func (t *Tool) Read(ctx context.Context, args RequestArgs) (*Content, error) {
	url := args.Url
	url = fmt.Sprintf("%s/%s", jinaHost, url)
	logger.FromContext(ctx).Debug("jina reader start", "url", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create jina reader request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send jina reader request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read jina reader response body: %w", err)
		}
		logger.FromContext(ctx).Error("jina reader read url error", "url", url, "status", resp.Status, "err", string(respData))
		return nil, fmt.Errorf("jina reader read URL '%s' Error: %s, %v", url, resp.Status, string(respData))
	}

	content := &Result{}

	if err := json.NewDecoder(resp.Body).Decode(content); err != nil {
		return nil, fmt.Errorf("failed to decode jina reader response: %w", err)
	}
	return &content.Data, nil
}
