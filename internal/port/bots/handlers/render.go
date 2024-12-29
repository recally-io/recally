package handlers

import (
	"bytes"
	"recally/internal/pkg/logger"
	"sync"

	tgmd "github.com/Mad-Pixels/goldmark-tgmd"
	"github.com/yuin/goldmark"
)

var mdParser goldmark.Markdown

var once sync.Once

func convertToTGMarkdown(md string) string {
	once.Do(func() {
		mdParser = tgmd.TGMD()
		tgmd.Config.UpdateHeading1(tgmd.Element{
			Prefix:  "\n📌 ",
			Postfix: "\n",
			Style:   tgmd.BoldTg,
		})
		tgmd.Config.UpdateHeading2(tgmd.Element{
			Prefix: "✏ ",
			Style:  tgmd.BoldTg,
		})
		tgmd.Config.UpdateHeading2(tgmd.Element{
			Prefix: "📚 ",
			Style:  tgmd.BoldTg,
		})
	})

	mdBytes := []byte(md)
	var buf bytes.Buffer
	if err := mdParser.Convert(mdBytes, &buf); err != nil {
		logger.Default.Warn("Failed to convert markdown to HTML", "error", err.Error())
		buf = *bytes.NewBuffer(mdBytes)
		return md
	}
	return buf.String()
}
