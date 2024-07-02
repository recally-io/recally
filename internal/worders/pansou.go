package worders

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
)

const pansearchHost = "https://www.pansearch.me"

func getBuildId() (string, error) {
	client := newHttpClient()
	resp, err := client.Get(pansearchHost)
	if err != nil {
		slog.Error("Failed to send request", "err", err, "url", pansearchHost)
		return "", err
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to fetch %s: %s, %v", pansearchHost, resp.Status, string(respData))
	}

	pattern := regexp.MustCompile(`"buildId":"(.*?)"`)
	matches := pattern.FindStringSubmatch(string(respData))
	if len(matches) < 2 {
		return "", fmt.Errorf("Failed to find buildId in response")
	}
	return matches[1], nil
}

type PansearchData struct {
	Content string  `json:"content"`
	Id      float32 `json:"id"`
	Image   string  `json:"image"`
	Pan     string  `json:"pan"`
	Time    string  `json:"time"`
}

func SearchAliPan(query string) ([]*PansearchData, error) {
	buildId, err := getBuildId()
	if err != nil {
		return nil, err
	}

	client := newHttpClient()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/_next/data/%s/search.json", pansearchHost, buildId), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("pan", "aliyundrive")
	q.Add("keyword", query)
	q.Add("offset", "0")
	q.Add("limit", "10")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	r := struct {
		PageProps struct {
			Data struct {
				Data  []*PansearchData `json:"data"`
				Total int              `json:"total"`
			} `json:"data"`
		} `json:"pageProps"`
	}{}

	// data, _ := io.ReadAll(resp.Body)
	// slog.Info("data", "data", string(data))
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return r.PageProps.Data.Data, nil
}
