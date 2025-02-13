package assistants

import (
	"context"
	"fmt"
	"recally/internal/core/queue"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/tools"
	"time"
)

type Service struct {
	llm   *llms.LLM
	dao   dao
	queue *queue.Queue
}

func NewService(llm *llms.LLM, queue *queue.Queue) *Service {
	return &Service{
		llm:   llm,
		dao:   db.New(),
		queue: queue,
	}
}

func (s *Service) ListModels(ctx context.Context) ([]llms.Model, error) {
	cacheKey := cache.NewCacheKey("list-models", "")
	if models, ok := cache.Get[[]llms.Model](ctx, cache.MemCache, cacheKey); ok {
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
