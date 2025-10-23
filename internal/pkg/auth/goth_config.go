package auth

import (
	"context"
	"fmt"

	"recally/internal/pkg/auth/adapter"
	"recally/internal/pkg/config"
	"recally/internal/pkg/db"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

// gothAdapterWrapper wraps the adapter.GothAdapter to match auth.OAuthAdapter interface
type gothAdapterWrapper struct {
	underlying *adapter.GothAdapter
}

func (w *gothAdapterWrapper) GetAuthURL(ctx context.Context, tx db.DBTX, provider string) (string, error) {
	return w.underlying.GetAuthURL(ctx, tx, provider)
}

func (w *gothAdapterWrapper) HandleCallback(ctx context.Context, tx db.DBTX, provider, code, state string) (any, error) {
	return w.underlying.HandleCallback(ctx, tx, provider, code, state)
}

func (w *gothAdapterWrapper) ListProviders() []string {
	return w.underlying.ListProviders()
}

// InitGothAdapter creates and configures a new Goth OAuth adapter
// with all enabled OAuth providers from the application configuration.
//
// This function:
// - Creates a new GothAdapter instance
// - Registers GitHub provider if configured
// - Registers Google provider if configured
// - Returns the configured adapter ready for use
//
// The adapter manages OAuth authentication flows including:
// - Generating secure auth URLs with CSRF protection
// - Handling OAuth callbacks
// - Validating state tokens
// - Fetching user information from providers
//
// Example usage:
//
//	adapter := InitGothAdapter(db.New())
//	authURL, err := adapter.GetAuthURL(ctx, "google")
func InitGothAdapter(dao *db.Queries) OAuthAdapter {
	gothAdapter := adapter.NewGothAdapter(dao)

	// Get service FQDN for callback URLs
	fqdn := config.Settings.Service.Fqdn

	// Register GitHub provider if configured
	if config.Settings.OAuths.Github.Key != "" && config.Settings.OAuths.Github.Secret != "" {
		githubProvider := github.New(
			config.Settings.OAuths.Github.Key,
			config.Settings.OAuths.Github.Secret,
			fmt.Sprintf("%s/api/v1/oauth/github/callback", fqdn),
			"user:email",
		)
		gothAdapter.RegisterProvider(githubProvider)
	}

	// Register Google provider if configured
	if config.Settings.OAuths.Google.Key != "" && config.Settings.OAuths.Google.Secret != "" {
		googleProvider := google.New(
			config.Settings.OAuths.Google.Key,
			config.Settings.OAuths.Google.Secret,
			fmt.Sprintf("%s/api/v1/oauth/google/callback", fqdn),
			"email", "profile",
		)
		gothAdapter.RegisterProvider(googleProvider)
	}

	// Clear any previously registered providers from goth.UseProviders
	// to ensure we only use providers registered in our adapter
	goth.ClearProviders()

	// Wrap the adapter to match the OAuthAdapter interface
	return &gothAdapterWrapper{underlying: gothAdapter}
}
