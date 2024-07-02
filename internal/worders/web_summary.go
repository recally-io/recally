package worders

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

const elmoSummaryUrl = "https://www.elmo.chat/api/v1/prompt"

func WebSummary(url string) (string, error) {
	return elmoSummary(url, "")
}

func elmoSummary(url, pageContent string) (string, error) {
	if pageContent == "" {
		content, err := WebReader(url)
		if err != nil {
			return "", fmt.Errorf("failed to get content: %w", err)
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
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

	client := newHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respData, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %w", err)
		}
		slog.Error("summary url error", "url", url, "status", resp.Status, "err", string(respData))
		return "", fmt.Errorf("failed to get summary: %s, %v", resp.Status, string(respData))
	}

	return processStreamResponse(resp.Body)
}

func processStreamResponse(body io.ReadCloser) (string, error) {
	var summary string
	reader := bufio.NewReader(body)
	defer body.Close()
	for {
		// Read until the next newline
		line, err := reader.ReadBytes('\n')
		if err != nil {
			// If the error is EOF, send it to errChan and return
			if errors.Is(err, io.EOF) {
				summary = strings.ReplaceAll(summary, "\\n", "\n")
				return summary, nil
			}
			return "", err
		}
		if bytes.HasPrefix(line, []byte("0:")) {
			// using regex to extract the summary text, origin text will be '0:"som text"', and I only want 'some text'
			summary += string(line[3 : len(line)-2])
		}
	}
}
