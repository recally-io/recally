package auth

import (
	"context"
	"fmt"
	"recally/internal/pkg/contexts"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
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

func LoadUserByID(ctx context.Context, tx db.DBTX, userID uuid.UUID) (*UserDTO, error) {
	return New().GetUserById(ctx, tx, userID)
}

func LoadUser(ctx context.Context, tx db.DBTX, userID uuid.UUID) (*UserDTO, error) {
	user, err := LoadUserFromContext(ctx)
	if err != nil {
		user, err = LoadUserByID(ctx, tx, userID)
	}
	return user, err
}
