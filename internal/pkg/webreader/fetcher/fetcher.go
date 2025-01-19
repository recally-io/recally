package fetcher

import "fmt"

type FecherType string

const (
	TypeHttp       FecherType = "http"
	TypeJinaReader FecherType = "jinaReader"
	TypeBrowser    FecherType = "browser"
	TypeNil        FecherType = ""
)

type FetchOptions struct {
	FecherType   FecherType `json:"fetcher_type"`   // the type of fetcher
	IsProxyImage bool       `json:"is_proxy_image"` // if true, the image will be proxied
	Force        bool       `json:"force"`          // if true, the fetcher will re-fetch the content
	IgnoreCache  bool       `json:"ignore_cache"`   // if true, the fetcher will ignore cache
}

func (o *FetchOptions) String() string {
	return fmt.Sprintf("%s-%t", o.FecherType, o.IsProxyImage)
}
