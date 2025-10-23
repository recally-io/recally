package adapter

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/markbates/goth"

	"recally/internal/pkg/db"
)

type GothAdapter struct {
	dao       *db.Queries
	providers map[string]goth.Provider
}

func NewGothAdapter(dao *db.Queries) *GothAdapter {
	return &GothAdapter{
		dao:       dao,
		providers: make(map[string]goth.Provider),
	}
}

func (a *GothAdapter) RegisterProvider(provider goth.Provider) {
	a.providers[provider.Name()] = provider
}

func (a *GothAdapter) GetAuthURL(ctx context.Context, tx db.DBTX, provider string) (string, error) {
	p, ok := a.providers[provider]
	if !ok {
		return "", fmt.Errorf("provider %s not found", provider)
	}

	// Generate secure state
	state, err := generateSecureState()
	if err != nil {
		return "", fmt.Errorf("failed to generate secure state: %w", err)
	}

	// Store state in database with 5-minute expiration
	expiresAt := time.Now().Add(5 * time.Minute)
	err = a.dao.CreateOAuthState(ctx, tx, db.CreateOAuthStateParams{
		State:       state,
		Provider:    provider,
		RedirectUrl: pgtype.Text{}, // Optional redirect URL
		ExpiresAt: pgtype.Timestamp{
			Time:  expiresAt,
			Valid: true,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to store OAuth state: %w", err)
	}

	// Get auth URL from provider
	sess, err := p.BeginAuth(state)
	if err != nil {
		return "", fmt.Errorf("failed to begin auth: %w", err)
	}

	authURL, err := sess.GetAuthURL()
	if err != nil {
		return "", fmt.Errorf("failed to get auth URL: %w", err)
	}

	return authURL, nil
}

func (a *GothAdapter) HandleCallback(ctx context.Context, tx db.DBTX, provider, code, state string) (OAuth2User, error) {
	// Validate state (CSRF protection)
	storedState, err := a.dao.GetOAuthState(ctx, tx, state)
	if err != nil {
		return OAuth2User{}, fmt.Errorf("invalid state: %w", err)
	}

	if storedState.Provider != provider {
		return OAuth2User{}, fmt.Errorf("provider mismatch: expected %s, got %s", storedState.Provider, provider)
	}

	// Check if state has expired
	if storedState.ExpiresAt.Valid && time.Now().After(storedState.ExpiresAt.Time) {
		return OAuth2User{}, fmt.Errorf("state expired")
	}

	// Delete used state (one-time use)
	_ = a.dao.DeleteOAuthState(ctx, tx, state)

	// Get provider and complete auth
	p, ok := a.providers[provider]
	if !ok {
		return OAuth2User{}, fmt.Errorf("provider %s not found", provider)
	}

	// Authorize with code
	sess, err := p.BeginAuth(state)
	if err != nil {
		return OAuth2User{}, fmt.Errorf("failed to begin auth: %w", err)
	}

	params := url.Values{}
	params.Set("code", code)
	_, err = sess.Authorize(p, params)
	if err != nil {
		return OAuth2User{}, fmt.Errorf("failed to authorize: %w", err)
	}

	gothUser, err := p.FetchUser(sess)
	if err != nil {
		return OAuth2User{}, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Marshal raw data
	rawData, _ := json.Marshal(gothUser.RawData)

	// Convert to our OAuth2User type
	return OAuth2User{
		Provider:       gothUser.Provider,
		ID:             gothUser.UserID,
		Name:           gothUser.Name,
		Email:          gothUser.Email,
		Avatar:         gothUser.AvatarURL,
		AccessToken:    gothUser.AccessToken,
		RefreshToken:   gothUser.RefreshToken,
		TokenExpiresAt: gothUser.ExpiresAt,
		RawData:        rawData,
	}, nil
}

func (a *GothAdapter) ListProviders() []string {
	names := make([]string, 0, len(a.providers))
	for name := range a.providers {
		names = append(names, name)
	}
	return names
}

// generateSecureState generates a cryptographically secure random state string
// Uses crypto/rand (NOT math/rand) for security
func generateSecureState() (string, error) {
	b := make([]byte, 32) // 32 bytes = 256 bits of entropy
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
