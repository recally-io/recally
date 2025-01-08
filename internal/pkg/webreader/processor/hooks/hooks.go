package hooks

import (
	"recally/internal/pkg/config"
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
	hooks := []md.BeforeHook{}

	if config.Settings.S3.Enabled {
		imageHook := NewImageHook(host)
		hooks = append(hooks, imageHook.Process)
	}

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

func GetReadabilityHooks(host string) []ReadabilityHook {
	hooks := []ReadabilityHook{}
	domainHooks, ok := readabilityHooks[host]
	if !ok {
		return hooks
	}
	return append(domainHooks, hooks...)
}
