package handlers

import (
	"context"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	GetUser(ctx context.Context, userID string) (*User, error)
	CreateUser(ctx context.Context, userName string, userID string) (*User, error)
	UpdateUser(ctx context.Context, user User) (*User, error)
}

type repository struct {
	db *db.Queries
}

func NewRepository(pool *db.Pool) Repository {
	return &repository{db: db.New(pool)}
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
