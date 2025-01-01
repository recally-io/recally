package bookmarks

import (
	"fmt"
)

// Common errors
var (
	ErrNotFound     = fmt.Errorf("bookmark not found")
	ErrDuplicate    = fmt.Errorf("bookmark already exists")
	ErrInvalidInput = fmt.Errorf("invalid input")
	ErrUnauthorized = fmt.Errorf("unauthorized access")
)

type FecherType string

const (
	HttpFetcher    FecherType = "http"
	JinaFetcher    FecherType = "jina"
	BrowserFetcher FecherType = "browser"
)
