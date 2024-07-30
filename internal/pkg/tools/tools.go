package tools

import (
	"context"

	"github.com/invopop/jsonschema"
	"github.com/sashabaranov/go-openai"
)

type Tool interface {
	Invoke(ctx context.Context, args string) (string, error)
	Schema() *openai.FunctionDefinition
	LLMName() string
}

type BaseTool struct {
	Name        string
	Description string
	Parameters  any
}

func (t *BaseTool) LLMName() string {
	return t.Name
}

func (t *BaseTool) Schema() *openai.FunctionDefinition {
	r := new(jsonschema.Reflector)
	schema := r.Reflect(t.Parameters)
	var paramsSchema *jsonschema.Schema

	if len(schema.Definitions) > 0 {
		for _, def := range schema.Definitions {
			paramsSchema = def
			break
		}
	} else {
		paramsSchema = schema
	}
	return &openai.FunctionDefinition{
		Name:        t.Name,
		Description: t.Description,
		Parameters:  paramsSchema,
	}
}
