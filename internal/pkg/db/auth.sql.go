// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: auth.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createOAuthConnection = `-- name: CreateOAuthConnection :one
INSERT INTO auth_user_oauth_connections (
    user_id, provider, provider_user_id, provider_email, 
    access_token, refresh_token, token_expires_at, provider_data
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, user_id, provider, provider_user_id, provider_email, access_token, refresh_token, token_expires_at, provider_data, created_at, updated_at
`

type CreateOAuthConnectionParams struct {
	UserID         uuid.UUID
	Provider       string
	ProviderUserID string
	ProviderEmail  pgtype.Text
	AccessToken    pgtype.Text
	RefreshToken   pgtype.Text
	TokenExpiresAt pgtype.Timestamptz
	ProviderData   []byte
}

func (q *Queries) CreateOAuthConnection(ctx context.Context, db DBTX, arg CreateOAuthConnectionParams) (AuthUserOauthConnection, error) {
	row := db.QueryRow(ctx, createOAuthConnection,
		arg.UserID,
		arg.Provider,
		arg.ProviderUserID,
		arg.ProviderEmail,
		arg.AccessToken,
		arg.RefreshToken,
		arg.TokenExpiresAt,
		arg.ProviderData,
	)
	var i AuthUserOauthConnection
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Provider,
		&i.ProviderUserID,
		&i.ProviderEmail,
		&i.AccessToken,
		&i.RefreshToken,
		&i.TokenExpiresAt,
		&i.ProviderData,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, email, phone, password_hash, activate_assistant_id, activate_thread_id, status, settings)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, uuid, username, password_hash, email, activate_assistant_id, activate_thread_id, status, created_at, updated_at, phone, settings
`

type CreateUserParams struct {
	Username            pgtype.Text
	Email               pgtype.Text
	Phone               pgtype.Text
	PasswordHash        pgtype.Text
	ActivateAssistantID pgtype.UUID
	ActivateThreadID    pgtype.UUID
	Status              string
	Settings            []byte
}

func (q *Queries) CreateUser(ctx context.Context, db DBTX, arg CreateUserParams) (User, error) {
	row := db.QueryRow(ctx, createUser,
		arg.Username,
		arg.Email,
		arg.Phone,
		arg.PasswordHash,
		arg.ActivateAssistantID,
		arg.ActivateThreadID,
		arg.Status,
		arg.Settings,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Username,
		&i.PasswordHash,
		&i.Email,
		&i.ActivateAssistantID,
		&i.ActivateThreadID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Phone,
		&i.Settings,
	)
	return i, err
}

const deleteExpiredRevokedTokens = `-- name: DeleteExpiredRevokedTokens :exec
DELETE FROM auth_revoked_tokens 
WHERE expires_at < CURRENT_TIMESTAMP
`

func (q *Queries) DeleteExpiredRevokedTokens(ctx context.Context, db DBTX) error {
	_, err := db.Exec(ctx, deleteExpiredRevokedTokens)
	return err
}

const deleteOAuthConnection = `-- name: DeleteOAuthConnection :exec
DELETE FROM auth_user_oauth_connections 
WHERE user_id = $1 AND provider = $2
`

type DeleteOAuthConnectionParams struct {
	UserID   uuid.UUID
	Provider string
}

func (q *Queries) DeleteOAuthConnection(ctx context.Context, db DBTX, arg DeleteOAuthConnectionParams) error {
	_, err := db.Exec(ctx, deleteOAuthConnection, arg.UserID, arg.Provider)
	return err
}

const deleteRevokedToken = `-- name: DeleteRevokedToken :exec
DELETE FROM auth_revoked_tokens 
WHERE jti = $1 AND user_id = $2
`

type DeleteRevokedTokenParams struct {
	Jti    uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteRevokedToken(ctx context.Context, db DBTX, arg DeleteRevokedTokenParams) error {
	_, err := db.Exec(ctx, deleteRevokedToken, arg.Jti, arg.UserID)
	return err
}

const deleteUserById = `-- name: DeleteUserById :exec
DELETE FROM users WHERE uuid = $1
`

func (q *Queries) DeleteUserById(ctx context.Context, db DBTX, argUuid uuid.UUID) error {
	_, err := db.Exec(ctx, deleteUserById, argUuid)
	return err
}

const getOAuthConnectionByProviderAndProviderID = `-- name: GetOAuthConnectionByProviderAndProviderID :one
SELECT id, user_id, provider, provider_user_id, provider_email, access_token, refresh_token, token_expires_at, provider_data, created_at, updated_at FROM auth_user_oauth_connections 
WHERE provider = $1 AND provider_user_id = $2
`

type GetOAuthConnectionByProviderAndProviderIDParams struct {
	Provider       string
	ProviderUserID string
}

func (q *Queries) GetOAuthConnectionByProviderAndProviderID(ctx context.Context, db DBTX, arg GetOAuthConnectionByProviderAndProviderIDParams) (AuthUserOauthConnection, error) {
	row := db.QueryRow(ctx, getOAuthConnectionByProviderAndProviderID, arg.Provider, arg.ProviderUserID)
	var i AuthUserOauthConnection
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Provider,
		&i.ProviderUserID,
		&i.ProviderEmail,
		&i.AccessToken,
		&i.RefreshToken,
		&i.TokenExpiresAt,
		&i.ProviderData,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getOAuthConnectionByUserAndProvider = `-- name: GetOAuthConnectionByUserAndProvider :one
SELECT id, user_id, provider, provider_user_id, provider_email, access_token, refresh_token, token_expires_at, provider_data, created_at, updated_at FROM auth_user_oauth_connections 
WHERE user_id = $1 AND provider = $2
`

type GetOAuthConnectionByUserAndProviderParams struct {
	UserID   uuid.UUID
	Provider string
}

func (q *Queries) GetOAuthConnectionByUserAndProvider(ctx context.Context, db DBTX, arg GetOAuthConnectionByUserAndProviderParams) (AuthUserOauthConnection, error) {
	row := db.QueryRow(ctx, getOAuthConnectionByUserAndProvider, arg.UserID, arg.Provider)
	var i AuthUserOauthConnection
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Provider,
		&i.ProviderUserID,
		&i.ProviderEmail,
		&i.AccessToken,
		&i.RefreshToken,
		&i.TokenExpiresAt,
		&i.ProviderData,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, uuid, username, password_hash, email, activate_assistant_id, activate_thread_id, status, created_at, updated_at, phone, settings FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, db DBTX, email pgtype.Text) (User, error) {
	row := db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Username,
		&i.PasswordHash,
		&i.Email,
		&i.ActivateAssistantID,
		&i.ActivateThreadID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Phone,
		&i.Settings,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, uuid, username, password_hash, email, activate_assistant_id, activate_thread_id, status, created_at, updated_at, phone, settings FROM users WHERE uuid = $1
`

func (q *Queries) GetUserById(ctx context.Context, db DBTX, argUuid uuid.UUID) (User, error) {
	row := db.QueryRow(ctx, getUserById, argUuid)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Username,
		&i.PasswordHash,
		&i.Email,
		&i.ActivateAssistantID,
		&i.ActivateThreadID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Phone,
		&i.Settings,
	)
	return i, err
}

const getUserByOAuthProviderId = `-- name: GetUserByOAuthProviderId :one
SELECT id, uuid, username, password_hash, email, activate_assistant_id, activate_thread_id, status, created_at, updated_at, phone, settings FROM users 
WHERE uuid = (SELECT user_id FROM auth_user_oauth_connections WHERE provider = $1 AND provider_user_id = $2)
`

type GetUserByOAuthProviderIdParams struct {
	Provider       string
	ProviderUserID string
}

func (q *Queries) GetUserByOAuthProviderId(ctx context.Context, db DBTX, arg GetUserByOAuthProviderIdParams) (User, error) {
	row := db.QueryRow(ctx, getUserByOAuthProviderId, arg.Provider, arg.ProviderUserID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Username,
		&i.PasswordHash,
		&i.Email,
		&i.ActivateAssistantID,
		&i.ActivateThreadID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Phone,
		&i.Settings,
	)
	return i, err
}

const getUserByPhone = `-- name: GetUserByPhone :one
SELECT id, uuid, username, password_hash, email, activate_assistant_id, activate_thread_id, status, created_at, updated_at, phone, settings FROM users WHERE phone = $1
`

func (q *Queries) GetUserByPhone(ctx context.Context, db DBTX, phone pgtype.Text) (User, error) {
	row := db.QueryRow(ctx, getUserByPhone, phone)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Username,
		&i.PasswordHash,
		&i.Email,
		&i.ActivateAssistantID,
		&i.ActivateThreadID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Phone,
		&i.Settings,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, uuid, username, password_hash, email, activate_assistant_id, activate_thread_id, status, created_at, updated_at, phone, settings FROM users WHERE username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, db DBTX, username pgtype.Text) (User, error) {
	row := db.QueryRow(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Username,
		&i.PasswordHash,
		&i.Email,
		&i.ActivateAssistantID,
		&i.ActivateThreadID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Phone,
		&i.Settings,
	)
	return i, err
}

const isTokenRevoked = `-- name: IsTokenRevoked :one
SELECT EXISTS (
    SELECT 1 FROM auth_revoked_tokens 
    WHERE jti = $1 AND user_id = $2
) AS is_revoked
`

type IsTokenRevokedParams struct {
	Jti    uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) IsTokenRevoked(ctx context.Context, db DBTX, arg IsTokenRevokedParams) (bool, error) {
	row := db.QueryRow(ctx, isTokenRevoked, arg.Jti, arg.UserID)
	var is_revoked bool
	err := row.Scan(&is_revoked)
	return is_revoked, err
}

const listOAuthConnectionsByUser = `-- name: ListOAuthConnectionsByUser :many
SELECT id, user_id, provider, provider_user_id, provider_email, access_token, refresh_token, token_expires_at, provider_data, created_at, updated_at FROM auth_user_oauth_connections 
WHERE user_id = $1 
ORDER BY created_at DESC
`

func (q *Queries) ListOAuthConnectionsByUser(ctx context.Context, db DBTX, userID uuid.UUID) ([]AuthUserOauthConnection, error) {
	rows, err := db.Query(ctx, listOAuthConnectionsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AuthUserOauthConnection
	for rows.Next() {
		var i AuthUserOauthConnection
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Provider,
			&i.ProviderUserID,
			&i.ProviderEmail,
			&i.AccessToken,
			&i.RefreshToken,
			&i.TokenExpiresAt,
			&i.ProviderData,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRevokedTokensByUser = `-- name: ListRevokedTokensByUser :many
SELECT jti, user_id, expires_at, revoked_at, reason FROM auth_revoked_tokens 
WHERE user_id = $1 
ORDER BY revoked_at DESC
`

func (q *Queries) ListRevokedTokensByUser(ctx context.Context, db DBTX, userID uuid.UUID) ([]AuthRevokedToken, error) {
	rows, err := db.Query(ctx, listRevokedTokensByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AuthRevokedToken
	for rows.Next() {
		var i AuthRevokedToken
		if err := rows.Scan(
			&i.Jti,
			&i.UserID,
			&i.ExpiresAt,
			&i.RevokedAt,
			&i.Reason,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsers = `-- name: ListUsers :many
SELECT id, uuid, username, password_hash, email, activate_assistant_id, activate_thread_id, status, created_at, updated_at, phone, settings FROM users ORDER BY created_at DESC
`

func (q *Queries) ListUsers(ctx context.Context, db DBTX) ([]User, error) {
	rows, err := db.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.Username,
			&i.PasswordHash,
			&i.Email,
			&i.ActivateAssistantID,
			&i.ActivateThreadID,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Phone,
			&i.Settings,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsersByStatus = `-- name: ListUsersByStatus :many
SELECT id, uuid, username, password_hash, email, activate_assistant_id, activate_thread_id, status, created_at, updated_at, phone, settings FROM users WHERE status = $1 ORDER BY created_at DESC
`

func (q *Queries) ListUsersByStatus(ctx context.Context, db DBTX, status string) ([]User, error) {
	rows, err := db.Query(ctx, listUsersByStatus, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.Username,
			&i.PasswordHash,
			&i.Email,
			&i.ActivateAssistantID,
			&i.ActivateThreadID,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Phone,
			&i.Settings,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const revokeToken = `-- name: RevokeToken :one
INSERT INTO auth_revoked_tokens (
    jti, user_id, expires_at, reason
) VALUES ($1, $2, $3, $4)
RETURNING jti, user_id, expires_at, revoked_at, reason
`

type RevokeTokenParams struct {
	Jti       uuid.UUID
	UserID    uuid.UUID
	ExpiresAt pgtype.Timestamptz
	Reason    pgtype.Text
}

func (q *Queries) RevokeToken(ctx context.Context, db DBTX, arg RevokeTokenParams) (AuthRevokedToken, error) {
	row := db.QueryRow(ctx, revokeToken,
		arg.Jti,
		arg.UserID,
		arg.ExpiresAt,
		arg.Reason,
	)
	var i AuthRevokedToken
	err := row.Scan(
		&i.Jti,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
		&i.Reason,
	)
	return i, err
}

const updateOAuthConnection = `-- name: UpdateOAuthConnection :one
UPDATE auth_user_oauth_connections SET 
    provider_email = $3,
    access_token = $4,
    refresh_token = $5,
    token_expires_at = $6,
    provider_data = $7,
    user_id = $8
WHERE provider_user_id = $1 AND provider = $2
RETURNING id, user_id, provider, provider_user_id, provider_email, access_token, refresh_token, token_expires_at, provider_data, created_at, updated_at
`

type UpdateOAuthConnectionParams struct {
	ProviderUserID string
	Provider       string
	ProviderEmail  pgtype.Text
	AccessToken    pgtype.Text
	RefreshToken   pgtype.Text
	TokenExpiresAt pgtype.Timestamptz
	ProviderData   []byte
	UserID         uuid.UUID
}

func (q *Queries) UpdateOAuthConnection(ctx context.Context, db DBTX, arg UpdateOAuthConnectionParams) (AuthUserOauthConnection, error) {
	row := db.QueryRow(ctx, updateOAuthConnection,
		arg.ProviderUserID,
		arg.Provider,
		arg.ProviderEmail,
		arg.AccessToken,
		arg.RefreshToken,
		arg.TokenExpiresAt,
		arg.ProviderData,
		arg.UserID,
	)
	var i AuthUserOauthConnection
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Provider,
		&i.ProviderUserID,
		&i.ProviderEmail,
		&i.AccessToken,
		&i.RefreshToken,
		&i.TokenExpiresAt,
		&i.ProviderData,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserById = `-- name: UpdateUserById :one
UPDATE users SET username = $2, email = $3, phone = $4, password_hash = $5,
  activate_assistant_id=$6, activate_thread_id=$7, status = $8, settings = $9
WHERE uuid = $1
RETURNING id, uuid, username, password_hash, email, activate_assistant_id, activate_thread_id, status, created_at, updated_at, phone, settings
`

type UpdateUserByIdParams struct {
	Uuid                uuid.UUID
	Username            pgtype.Text
	Email               pgtype.Text
	Phone               pgtype.Text
	PasswordHash        pgtype.Text
	ActivateAssistantID pgtype.UUID
	ActivateThreadID    pgtype.UUID
	Status              string
	Settings            []byte
}

func (q *Queries) UpdateUserById(ctx context.Context, db DBTX, arg UpdateUserByIdParams) (User, error) {
	row := db.QueryRow(ctx, updateUserById,
		arg.Uuid,
		arg.Username,
		arg.Email,
		arg.Phone,
		arg.PasswordHash,
		arg.ActivateAssistantID,
		arg.ActivateThreadID,
		arg.Status,
		arg.Settings,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Username,
		&i.PasswordHash,
		&i.Email,
		&i.ActivateAssistantID,
		&i.ActivateThreadID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Phone,
		&i.Settings,
	)
	return i, err
}
