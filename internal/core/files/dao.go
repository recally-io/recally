package files

import (
	"context"

	"recally/internal/pkg/db"

	"github.com/google/uuid"
)

// DAO provides data access operations for files.
type DAO interface {
	CreateFile(ctx context.Context, tx db.DBTX, arg db.CreateFileParams) (db.File, error)

	DeleteFile(ctx context.Context, tx db.DBTX, id uuid.UUID) error
	DeleteFileByOriginalURL(ctx context.Context, db db.DBTX, arg db.DeleteFileByOriginalURLParams) error
	DeleteFileByS3Key(ctx context.Context, tx db.DBTX, arg db.DeleteFileByS3KeyParams) error

	GetFileByID(ctx context.Context, tx db.DBTX, id uuid.UUID) (db.File, error)
	GetFileByOriginalURL(ctx context.Context, db db.DBTX, arg db.GetFileByOriginalURLParams) (db.File, error)
	GetFileByS3Key(ctx context.Context, db db.DBTX, arg db.GetFileByS3KeyParams) (db.File, error)
}
