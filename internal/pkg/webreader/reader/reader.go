package reader

import (
	"fmt"
	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/fetcher"
	"recally/internal/pkg/webreader/processor"
	"recally/internal/pkg/webreader/processor/hooks"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

func New(host string, opts fetcher.FetchOptions) (*webreader.Reader, error) {
	var readerFetcher webreader.Fetcher
	var err error
	switch opts.FecherType {
	case fetcher.TypeHttp:
		readerFetcher, err = fetcher.NewHTTPFetcher()
	case fetcher.TypeJinaReader:
		readerFetcher, err = fetcher.NewJinaFetcher()
	case fetcher.TypeBrowser:
		readerFetcher, err = fetcher.NewBrowserFetcher()
	}

	if err != nil {
		return nil, err
	}
	if readerFetcher == nil {
		return nil, fmt.Errorf("fetcher not found")
	}

	mdBeforeHooks := []md.BeforeHook{}
	if opts.IsProxyImage {
		imgHook := hooks.NewImageHook(host)
		mdBeforeHooks = append(mdBeforeHooks, imgHook.Process)
	}

	processors := []webreader.Processor{
		processor.NewReadabilityProcessor(),
		processor.NewMarkdownProcessor(host, processor.WithMarkdownBeforeHook(mdBeforeHooks...)),
	}

	return webreader.New(readerFetcher, processors...), nil
}
