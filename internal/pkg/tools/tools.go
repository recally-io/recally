package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/sashabaranov/go-openai"
)

type Tool interface {
	Invoke(ctx context.Context, args string) (string, error)
	Schema() *openai.FunctionDefinition
	LLMName() string
	LLMDescription() string
}

type BaseTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

func (t *BaseTool) LLMName() string {
	return t.Name
}

func (t *BaseTool) LLMDescription() string {
	return t.Description
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

// UnmarshalArgs unmarshals the args string into the params struct.
// params must be a pointer to a struct.
func (t *BaseTool) UnmarshalArgs(ctx context.Context, args string, params any) error {
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return fmt.Errorf("failed to unmarshal %s request: %w", t.Name, err)
	}
	return nil
}

func (t *BaseTool) MarshalResult(ctx context.Context, result any) (string, error) {
	b, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal %s result: %w", t.Name, err)
	}
	return string(b), nil
}
