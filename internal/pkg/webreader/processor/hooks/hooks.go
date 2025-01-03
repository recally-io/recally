package hooks

import (
	"recally/internal/pkg/webreader"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

type ReadabilityHook func(doc *goquery.Document, content *webreader.Content)

var (
	mdBeforeHooks    = map[string][]md.BeforeHook{}
	mdAfterHooks     = map[string][]md.Afterhook{}
	readabilityHooks = map[string][]ReadabilityHook{}
)

func init() {
	NewWeixinMP().RegisterHooks()
}

func GetMarkdownBeforeHooks(host string) []md.BeforeHook {
	hooks, ok := mdBeforeHooks[host]
	if !ok {
		return []md.BeforeHook{}
	}
	return hooks
}

func GetMarkdownAfterHooks(host string) []md.Afterhook {
	hooks, ok := mdAfterHooks[host]
	if !ok {
		return []md.Afterhook{}
	}
	return hooks
}

func GetReadabilityHooks(host string) []ReadabilityHook {
	hooks, ok := readabilityHooks[host]
	if !ok {
		return []ReadabilityHook{}
	}
	return hooks
}
