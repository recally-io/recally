package adapter

import (
	"context"
	"time"

	"recally/internal/pkg/db"
)

// OAuth2User represents user information from OAuth provider
type OAuth2User struct {
	Provider string `json:"provider"`
	ID       string `json:"id"`
	Name     string `json:"user_name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar_url"`

	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	TokenExpiresAt time.Time `json:"token_expires_at"`
	RawData        []byte    `json:"raw_data"`
}

// OAuthAdapter abstracts OAuth provider implementation
type OAuthAdapter interface {
	// GetAuthURL generates OAuth redirect URL with secure state
	// Requires a database transaction to store the state for CSRF protection
	GetAuthURL(ctx context.Context, tx db.DBTX, provider string) (url string, err error)

	// HandleCallback processes OAuth callback and returns user info
	// Requires a database transaction to validate state and prevent CSRF attacks
	HandleCallback(ctx context.Context, tx db.DBTX, provider, code, state string) (OAuth2User, error)

	// ListProviders returns available OAuth providers
	ListProviders() []string
}
