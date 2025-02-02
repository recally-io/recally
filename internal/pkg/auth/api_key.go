package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"recally/internal/pkg/db"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrInvalidApiKey = errors.New("401: invalid API key")

type ApiKeyDTO struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Name       string    `json:"name"`
	Prefix     string    `json:"prefix"`
	Hash       string    `json:"hash"`
	Scopes     []string  `json:"scopes"`
	ExpiresAt  time.Time `json:"expires_at"`
	LastUsedAt time.Time `json:"last_used_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (t *ApiKeyDTO) Load(d *db.AuthApiKey) {
	t.ID = d.ID
	t.UserID = d.UserID
	t.Name = d.Name
	t.Prefix = d.KeyPrefix
	t.Hash = d.KeyHash
	t.Scopes = d.Scopes
	if d.ExpiresAt.Valid {
		t.ExpiresAt = d.ExpiresAt.Time
	}
	if d.LastUsedAt.Valid {
		t.LastUsedAt = d.LastUsedAt.Time
	}
	t.CreatedAt = d.CreatedAt.Time
	t.UpdatedAt = d.UpdatedAt.Time
}

func (t *ApiKeyDTO) Dump() *db.AuthApiKey {
	return &db.AuthApiKey{
		ID:        t.ID,
		UserID:    t.UserID,
		Name:      t.Name,
		KeyPrefix: t.Prefix,
		KeyHash:   t.Hash,
		Scopes:    t.Scopes,
		ExpiresAt: pgtype.Timestamptz{
			Time:  t.ExpiresAt,
			Valid: t.ExpiresAt != time.Time{},
		},
		LastUsedAt: pgtype.Timestamptz{
			Time:  t.LastUsedAt,
			Valid: t.LastUsedAt != time.Time{},
		},
	}
}

func (s *Service) generateRandomApiKey(prefix string) string {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}

	// Encode to base64
	encoded := base64.URLEncoding.EncodeToString(bytes)

	// Remove any padding characters
	encoded = strings.TrimRight(encoded, "=")

	// Combine prefix with random part
	if prefix != "" {
		return prefix + "_" + encoded
	}
	return encoded
}

func (s *Service) CreateApiKey(ctx context.Context, tx db.DBTX, key *ApiKeyDTO) (*ApiKeyDTO, error) {
	if key.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("API key expiration time must be in the future")
	}

	if key.Hash == "" {
		// generate hash
		key.Hash = s.generateRandomApiKey(key.Prefix)
	}

	dbKey := key.Dump()
	params := db.CreateAPIKeyParams{
		UserID:    dbKey.UserID,
		Name:      dbKey.Name,
		KeyPrefix: dbKey.KeyPrefix,
		KeyHash:   fmt.Sprintf("ak-%s", dbKey.KeyHash),
		Scopes:    []string{"*"},
		ExpiresAt: dbKey.ExpiresAt,
	}

	apiKey, err := s.dao.CreateAPIKey(ctx, tx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	key.Load(&apiKey)
	return key, nil
}

func (s *Service) DeleteApiKey(ctx context.Context, tx db.DBTX, id uuid.UUID) error {
	if err := s.dao.DeleteAPIKey(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}
	return nil
}

func (s *Service) ListApiKeys(ctx context.Context, tx db.DBTX, prefix string, IsActive bool) ([]*ApiKeyDTO, error) {
	user, err := LoadUserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load user from context: %w", err)
	}

	apiKeys, err := s.dao.ListAPIKeys(ctx, tx, db.ListAPIKeysParams{
		UserID:   user.ID,
		Prefix:   pgtype.Text{String: prefix, Valid: prefix != ""},
		IsActive: pgtype.Bool{Bool: IsActive, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}

	result := make([]*ApiKeyDTO, len(apiKeys))
	for i, key := range apiKeys {
		dto := new(ApiKeyDTO)
		dto.Load(&key)
		result[i] = dto
	}
	return result, nil
}

func (s *Service) UpdateApiKeyLastUsed(ctx context.Context, tx db.DBTX, id uuid.UUID) error {
	if err := s.dao.UpdateAPIKeyLastUsed(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to update API key last used: %w", err)
	}
	return nil
}

func (s *Service) ValidateApiKey(ctx context.Context, tx db.DBTX, key string) (*UserDTO, error) {
	user, err := s.dao.GetUserByApiKey(ctx, tx, key)
	if err != nil {
		return nil, err
	}
	u := new(UserDTO)
	u.Load(&user)
	return u, nil
}
