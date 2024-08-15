package auth

import (
	"context"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
)

type dto interface {
	GetUserById(ctx context.Context, db db.DBTX, argUuid uuid.UUID) (db.User, error)
	CreateUser(ctx context.Context, db db.DBTX, arg db.CreateUserParams) (db.User, error)
}
