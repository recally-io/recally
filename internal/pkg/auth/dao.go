package auth

import (
	"context"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type dto interface {
	GetUserById(ctx context.Context, db db.DBTX, argUuid uuid.UUID) (db.User, error)
	CreateUser(ctx context.Context, db db.DBTX, arg db.CreateUserParams) (db.User, error)
	GetTelegramUser(ctx context.Context, db db.DBTX, telegram pgtype.Text) (db.User, error)
	InserUser(ctx context.Context, db db.DBTX, arg db.InserUserParams) (db.User, error)
	UpdateTelegramUser(ctx context.Context, db db.DBTX, arg db.UpdateTelegramUserParams) (db.User, error)
}
