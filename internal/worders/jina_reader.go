package worders

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type WebReaderContent struct {
	Url         string `json:"url"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
}

type jinaReaderResponse struct {
	Code   int              `json:"code"`
	Status float64          `json:"status"`
	Data   WebReaderContent `json:"data"`
}

type jinaSearcherResponse struct {
	Code   int                 `json:"code"`
	Status float64             `json:"status"`
	Data   []*WebReaderContent `json:"data"`
}

const (
	jinaReaderHost   = "https://r.jina.ai"
	jinaSearcherHost = "https://s.jina.ai"
)

func newHttpClient() *http.Client {
	return &http.Client{
		Timeout: 5 * 60 * time.Second,
	}
}

func WebReader(url string) (*WebReaderContent, error) {
	readerUrl := fmt.Sprintf("%s/%s", jinaReaderHost, url)
	req, err := http.NewRequest(http.MethodGet, readerUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := newHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		slog.Error("web reader read url error", "url", url, "status", resp.Status, "err", string(respData))
		return nil, fmt.Errorf("Read URL %s Error: %s, %v", readerUrl, resp.Status, string(respData))
	}

	content := &jinaReaderResponse{}

	if err := json.NewDecoder(resp.Body).Decode(content); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &content.Data, nil
}

func WebSearcher(query string) ([]*WebReaderContent, error) {
	searcherUrl := fmt.Sprintf("%s/%s", jinaSearcherHost, query)
	req, err := http.NewRequest(http.MethodGet, searcherUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := newHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		slog.Error("web searcher error", "query", query, "status", resp.Status, "err", string(respData))
		return nil, fmt.Errorf("Search web Error: %s, %v", resp.Status, string(respData))
	}

	content := &jinaSearcherResponse{}

	if err := json.NewDecoder(resp.Body).Decode(content); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return content.Data, nil
}
