package assistants

import (
	"context"
	"fmt"
	"time"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/llms"
	"vibrain/internal/pkg/tools"
)

type Service struct {
	llm *llms.LLM
	dao dao
}

func NewService(llm *llms.LLM) *Service {
	return &Service{
		llm: llm,
		dao: db.New(),
	}
}

func (s *Service) ListModels(ctx context.Context) ([]string, error) {
	cacheKey := cache.NewCacheKey("list-models", "")
	if models, ok := cache.Get[[]string](ctx, cache.MemCache, cacheKey); ok {
		return *models, nil
	}
	models, err := s.llm.ListModels(ctx)
	if err != nil {
		return nil, fmt.Errorf("list models error: %w", err)
	}
	cache.MemCache.Set(cacheKey, &models, time.Hour)
	return models, nil
}

func (s *Service) ListTools(ctx context.Context) ([]tools.BaseTool, error) {
	toolMappings := llms.AllToolMappings
	availableTools := make([]tools.BaseTool, 0, len(toolMappings))
	for _, tool := range toolMappings {
		availableTools = append(availableTools, tools.BaseTool{
			Name:        tool.LLMName(),
			Description: tool.LLMDescription(),
		})
	}
	return availableTools, nil
}
