package auth

import (
	"context"
	"fmt"
	"strings"
	"vibrain/internal/pkg/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func getOAuth2Config(provider string) (*oauth2.Config, error) {
	for _, oauth := range config.Settings.OAuths {
		if strings.EqualFold(oauth.Provider, provider) {
			return &oauth2.Config{
				ClientID:     oauth.Key,
				ClientSecret: oauth.Secret,
				Endpoint:     github.Endpoint,
				RedirectURL:  fmt.Sprintf("%s/oauth/github/callback", config.Settings.Service.Fqdn),
				Scopes:       []string{"user:email"},
			}, nil
		}
	}

	return nil, fmt.Errorf("oauth provider '%s' not found", provider)
}

func GetOAuth2RedirectURL(ctx context.Context, provider string) (string, error) {
	cfg, err := getOAuth2Config(provider)
	if err != nil {
		return "", fmt.Errorf("failed to get oauth config: %w", err)
	}
	authCodeUrl := cfg.AuthCodeURL("state:"+provider, oauth2.AccessTypeOnline)
	return authCodeUrl, nil
}

func GetOAuth2Token(ctx context.Context, provider, code string) (*oauth2.Token, error) {
	cfg, err := getOAuth2Config(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth config: %w", err)
	}
	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange oauth code: %w", err)
	}
	return token, nil
}
