package auth

import (
	"context"
	"fmt"
	"recally/internal/pkg/contexts"
)

func LoadUserFromContext(ctx context.Context) (*UserDTO, error) {
	user, ok := contexts.Get[*UserDTO](ctx, contexts.ContextKeyUser)
	if !ok {
		return nil, fmt.Errorf("failed to get user from context")
	}
	return user, nil
}

func SetUserToContext(ctx context.Context, user *UserDTO) context.Context {
	return context.WithValue(ctx, contexts.ContextKey(contexts.ContextKeyUser), user)
}
