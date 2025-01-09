package hooks

import (
	"recally/internal/pkg/webreader"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

type ReadabilityHook func(doc *goquery.Document, content *webreader.Content)

var (
	mdBeforeHooks          = map[string][]md.BeforeHook{}
	mdAfterHooks           = map[string][]md.Afterhook{}
	readabilityBeforeHooks = map[string][]ReadabilityHook{}
	readabilityAfterHooks  = map[string][]ReadabilityHook{}
)

func init() {
	NewWeixinMP().RegisterHooks()
}

func GetMarkdownBeforeHooks(host string) []md.BeforeHook {
	hooks := []md.BeforeHook{}
	domainHooks, ok := mdBeforeHooks[host]
	if !ok {
		return hooks
	}
	return append(domainHooks, hooks...)
}

func GetMarkdownAfterHooks(host string) []md.Afterhook {
	hooks := []md.Afterhook{}
	domainHooks, ok := mdAfterHooks[host]
	if !ok {
		return hooks
	}
	return append(domainHooks, hooks...)
}

func GetReadabilityBeforeHooks(host string) []ReadabilityHook {
	hooks := []ReadabilityHook{}
	domainHooks, ok := readabilityBeforeHooks[host]
	if !ok {
		return hooks
	}
	return append(domainHooks, hooks...)
}

func GetReadabilityAfterHooks(host string) []ReadabilityHook {
	hooks := []ReadabilityHook{}
	domainHooks, ok := readabilityAfterHooks[host]
	if !ok {
		return hooks
	}
	return append(domainHooks, hooks...)
}
