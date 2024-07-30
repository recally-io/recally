package llms

import (
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/tools"
	"vibrain/internal/pkg/tools/googlesearch"
	"vibrain/internal/pkg/tools/jinareader"
	"vibrain/internal/pkg/tools/jinasearcher"

	"github.com/sashabaranov/go-openai"
)

func DefaultLLMToolMappings() map[string]tools.Tool {
	allTools := []tools.Tool{
		jinareader.New(),
		jinasearcher.New(),
	}

	if config.Settings.GoogleSearch.ApiKey != "" && config.Settings.GoogleSearch.EngineID != "" {
		allTools = append(allTools, googlesearch.New(config.Settings.GoogleSearch.ApiKey, config.Settings.GoogleSearch.EngineID))
	}

	mappings := make(map[string]tools.Tool)
	for _, t := range allTools {
		mappings[t.LLMName()] = t
	}

	return mappings
}

func llmTools(mappings map[string]tools.Tool) []openai.Tool {
	toolList := make([]openai.Tool, 0)
	for _, t := range mappings {
		schema := t.Schema()
		toolList = append(toolList, openai.Tool{
			Type: "function",
			Function: &openai.FunctionDefinition{
				Name:        schema.Name,
				Description: schema.Description,
				Parameters:  schema.Parameters,
			},
		})
	}
	return toolList
}
