package auth

import (
	"context"
	"fmt"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/contexts"
	"recally/internal/pkg/db"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func LoadUserFromContext(ctx context.Context) (*UserDTO, error) {
	user, ok := contexts.Get[*UserDTO](ctx, contexts.ContextKeyUser)
	if !ok {
		return nil, fmt.Errorf("failed to get user from context")
	}
	return user, nil
}

func SetUserToContext(ctx context.Context, user *UserDTO) context.Context {
	ctx = contexts.Set(ctx, contexts.ContextKeyUser, user)
	ctx = contexts.Set(ctx, contexts.ContextKeyUserID, user.ID)
	ctx = contexts.Set(ctx, contexts.ContextKeyUserName, user.Username)
	return ctx
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

func LoadDummyUser() (*UserDTO, error) {
	ctx := context.Background()
	user, err := cache.RunInCache[UserDTO](ctx,
		cache.MemCache,
		cache.NewCacheKey("auth", "dummy_user"),
		time.Hour,
		func() (*UserDTO, error) {
			var user *UserDTO
			var err error
			if e := db.RunInTransaction(ctx, db.DefaultPool.Pool, func(ctx context.Context, tx pgx.Tx) error {
				user, err = New().GetDummyUser(ctx, tx)
				return err
			}); e != nil {
				return nil, e
			}
			return user, nil
		})

	return user, err
}

func GetContextWithDummyUser(ctx context.Context) (context.Context, error) {
	dummyUser, err := LoadDummyUser()
	if err != nil {
		return nil, err
	}
	return SetUserToContext(ctx, dummyUser), nil
}
