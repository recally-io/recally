package handlers

import (
	"context"
	"fmt"
	"strings"
	"vibrain/internal/pkg/contexts"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	GetUser(ctx context.Context, userID string) (*User, error)
	CreateUser(ctx context.Context, userName string, userID string) (*User, error)
	UpdateUser(ctx context.Context, user User) (*User, error)
	// GetOrCreateUser gets user by telegram userID, if user not found, creates new user with userName and userID
	GetOrCreateUser(ctx context.Context) (*User, error)
}

type repository struct {
	db *db.Queries
}

func NewRepository(db *db.Queries) Repository {
	return &repository{db: db}
}

func NewRepositoryFromContext(ctx context.Context) (Repository, error) {
	tx, ok := contexts.Get[pgx.Tx](ctx, contexts.ContextKeyTx)
	if ok {
		return NewRepository(db.New(tx)), nil
	}
	return nil, fmt.Errorf("failed to get db pool from context")
}

func (r *repository) GetUser(ctx context.Context, userID string) (*User, error) {
	user, err := r.db.GetTelegramUser(ctx, pgtype.Text{
		String: userID,
		Valid:  true,
	})
	if err != nil {
		return nil, err
	}
	return &User{
		ID:                  user.Uuid,
		Username:            user.Username.String,
		Telegram:            user.Telegram.String,
		ActivateAssistantID: user.ActivateAssistantID.Bytes,
		ActivateThreadID:    user.ActivateThreadID.Bytes,
	}, nil
}

func (r *repository) CreateUser(ctx context.Context, userName string, userID string) (*User, error) {
	params := db.InserUserParams{
		Username: pgtype.Text{String: userName, Valid: true},
		Telegram: pgtype.Text{String: userID, Valid: true},
	}
	user, err := r.db.InserUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:                  user.Uuid,
		Username:            user.Username.String,
		Telegram:            user.Telegram.String,
		ActivateAssistantID: user.ActivateAssistantID.Bytes,
	}, nil
}

func (r *repository) UpdateUser(ctx context.Context, user User) (*User, error) {
	dbUser, err := r.db.UpdateTelegramUser(ctx, db.UpdateTelegramUserParams{
		Telegram:            pgtype.Text{String: user.Telegram, Valid: true},
		ActivateAssistantID: pgtype.UUID{Bytes: user.ActivateAssistantID, Valid: user.ActivateAssistantID != uuid.Nil},
		ActivateThreadID:    pgtype.UUID{Bytes: user.ActivateThreadID, Valid: user.ActivateThreadID != uuid.Nil},
	})
	if err != nil {
		return nil, err
	}
	return &User{
		ID:                  dbUser.Uuid,
		Username:            dbUser.Username.String,
		Telegram:            dbUser.Telegram.String,
		ActivateAssistantID: dbUser.ActivateAssistantID.Bytes,
		ActivateThreadID:    dbUser.ActivateThreadID.Bytes,
	}, nil
}

func (r *repository) GetOrCreateUser(ctx context.Context) (*User, error) {
	userID, ok := contexts.Get[string](ctx, contexts.ContextKeyUserID)
	if !ok {
		return nil, fmt.Errorf("failed to get userID from context")
	}
	user, err := r.GetUser(ctx, userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			userName := ctx.Value(contexts.ContextKey(contexts.ContextKeyUserName)).(string)
			user, err = r.CreateUser(ctx, userName, userID)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
	}
	return user, nil
}
