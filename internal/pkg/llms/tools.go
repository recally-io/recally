package llms

import (
	"recally/internal/pkg/config"
	"recally/internal/pkg/tools"
	"recally/internal/pkg/tools/googlesearch"
	"recally/internal/pkg/tools/jinareader"
	"recally/internal/pkg/tools/jinasearcher"

	"github.com/sashabaranov/go-openai"
)

var AllToolMappings map[string]tools.Tool

func init() {
	AllToolMappings = defaultLLMToolMappings()
}

func defaultLLMToolMappings() map[string]tools.Tool {
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
