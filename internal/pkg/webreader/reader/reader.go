package reader

import (
	"fmt"
	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/fetcher"
	"recally/internal/pkg/webreader/processor"
)

func New(fetcherType fetcher.FecherType, host string) (*webreader.Reader, error) {
	var readerFetcher webreader.Fetcher
	var err error
	switch fetcherType {
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

	processors := []webreader.Processor{
		processor.NewReadabilityProcessor(),
		processor.NewMarkdownProcessor(host),
	}

	return webreader.New(readerFetcher, processors...), nil
}
