package reader

import (
	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/fetcher"
	"recally/internal/pkg/webreader/processor"
)

func New(fetcherType fetcher.FecherType, host string) (*webreader.Reader, error) {
	var readerFetcher webreader.Fetcher
	var err error
	processors := []webreader.Processor{
		processor.NewReadabilityProcessor(),
	}

	switch fetcherType {
	case fetcher.TypeHttp:
		readerFetcher, err = fetcher.NewHTTPFetcher()
		if err != nil {
			return nil, err
		}
		processors = append(processors, processor.NewMarkdownProcessor(host))
	case fetcher.TypeJinaReader:
		readerFetcher, err = fetcher.NewJinaFetcher()
		if err != nil {
			return nil, err
		}
	case fetcher.TypeBrowser:
		readerFetcher, err = fetcher.NewBrowserFetcher()
		if err != nil {
			return nil, err
		}
		processors = append(processors, processor.NewMarkdownProcessor(host))
	}

	return webreader.New(readerFetcher, processors...), nil
}
