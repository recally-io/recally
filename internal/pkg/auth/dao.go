package auth

import (
	"context"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type dto interface {
	CreateOAuthConnection(ctx context.Context, db db.DBTX, arg db.CreateOAuthConnectionParams) (db.AuthUserOauthConnection, error)
	CreateUser(ctx context.Context, db db.DBTX, arg db.CreateUserParams) (db.User, error)
	DeleteExpiredRevokedTokens(ctx context.Context, db db.DBTX) error
	DeleteOAuthConnection(ctx context.Context, db db.DBTX, arg db.DeleteOAuthConnectionParams) error
	DeleteRevokedToken(ctx context.Context, db db.DBTX, arg db.DeleteRevokedTokenParams) error
	DeleteUserById(ctx context.Context, db db.DBTX, argUuid uuid.UUID) error
	GetOAuthConnectionByUserAndProvider(ctx context.Context, db db.DBTX, arg db.GetOAuthConnectionByUserAndProviderParams) (db.AuthUserOauthConnection, error)
	GetOAuthConnectionByProviderAndProviderID(ctx context.Context, db db.DBTX, arg db.GetOAuthConnectionByProviderAndProviderIDParams) (db.AuthUserOauthConnection, error)
	GetUserByEmail(ctx context.Context, db db.DBTX, email pgtype.Text) (db.User, error)
	GetUserById(ctx context.Context, db db.DBTX, argUuid uuid.UUID) (db.User, error)
	GetUserByOAuthProviderId(ctx context.Context, db db.DBTX, arg db.GetUserByOAuthProviderIdParams) (db.User, error)
	GetUserByPhone(ctx context.Context, db db.DBTX, phone pgtype.Text) (db.User, error)
	GetUserByUsername(ctx context.Context, db db.DBTX, username pgtype.Text) (db.User, error)
	IsTokenRevoked(ctx context.Context, db db.DBTX, arg db.IsTokenRevokedParams) (bool, error)
	ListOAuthConnectionsByUser(ctx context.Context, db db.DBTX, userID uuid.UUID) ([]db.AuthUserOauthConnection, error)
	ListRevokedTokensByUser(ctx context.Context, db db.DBTX, userID uuid.UUID) ([]db.AuthRevokedToken, error)
	ListUsers(ctx context.Context, db db.DBTX) ([]db.User, error)
	ListUsersByStatus(ctx context.Context, db db.DBTX, status string) ([]db.User, error)
	RevokeToken(ctx context.Context, db db.DBTX, arg db.RevokeTokenParams) (db.AuthRevokedToken, error)
	UpdateOAuthConnection(ctx context.Context, db db.DBTX, arg db.UpdateOAuthConnectionParams) (db.AuthUserOauthConnection, error)
	UpdateUserById(ctx context.Context, db db.DBTX, arg db.UpdateUserByIdParams) (db.User, error)

	OwnerTransferBookmark(ctx context.Context, db db.DBTX, arg db.OwnerTransferBookmarkParams) error

	CreateAPIKey(ctx context.Context, db db.DBTX, arg db.CreateAPIKeyParams) (db.AuthApiKey, error)
	DeleteAPIKey(ctx context.Context, db db.DBTX, id uuid.UUID) error
	ListAPIKeys(ctx context.Context, db db.DBTX, arg db.ListAPIKeysParams) ([]db.AuthApiKey, error)
	UpdateAPIKeyLastUsed(ctx context.Context, db db.DBTX, id uuid.UUID) error
	GetUserByApiKey(ctx context.Context, db db.DBTX, keyHash string) (db.User, error)
}
