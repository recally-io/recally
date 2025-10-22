package googlesearch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"recally/internal/pkg/tools"
)

type Tool struct {
	tools.BaseTool
	apiKey         string
	searchEngineID string
	client         *http.Client
}

type Option func(*Tool)

func WithHttpClient(client *http.Client) Option {
	return func(gs *Tool) {
		gs.client = client
	}
}

func New(apiKey, searchEngineID string) *Tool {
	return &Tool{
		apiKey:         apiKey,
		searchEngineID: searchEngineID,
		client:         http.DefaultClient,
		BaseTool: tools.BaseTool{
			Name:        "googlesearch",
			Description: "A tool that using google search API to search for information.",
			Parameters:  &RequestArgs{},
		},
	}
}

type RequestArgs struct {
	Q string `json:"q" jsonschema_description:"The search query."`
}

type ResultItem struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}

type Result struct {
	Items []ResultItem `json:"items"`
	Kind  string       `json:"kind"`
}

func (t *Tool) Invoke(ctx context.Context, args string) (string, error) {
	var params RequestArgs
	if err := t.UnmarshalArgs(ctx, args, &params); err != nil {
		return "", err
	}

	result, err := t.Search(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to invoke tool: %w", err)
	}

	return t.MarshalResult(ctx, result)
}

func (t *Tool) Search(ctx context.Context, args RequestArgs) (*Result, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/customsearch/v1", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create google search request: %w", err)
	}

	params := req.URL.Query()
	params.Add("key", t.apiKey)
	params.Add("cx", t.searchEngineID)
	params.Add("q", args.Q)
	req.URL.RawQuery = params.Encode()

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send google search request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	result := new(Result)
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("failed to get google search response: %w", err)
	}

	return result, nil
}
