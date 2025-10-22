package auth

import (
	"context"
	"fmt"
	"time"

	"recally/internal/pkg/config"
	"recally/internal/pkg/logger"

	"golang.org/x/oauth2"
)

type ProviderType string

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

type OAuthProvider interface {
	GetConfig() *oauth2.Config
	GetRedirectURL() string
	GetToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetUser(ctx context.Context, token *oauth2.Token) (OAuth2User, error)
}

func GetOAuthProvider(name string) (OAuthProvider, error) {
	if name == "github" {
		provider := NewOAuthProviderGithub(config.Settings.OAuths.Github)

		return provider, nil
	}

	return nil, fmt.Errorf("oauth provider '%s' not found", name)
}

type oAuthProvider struct {
	Name          ProviderType
	cfg           config.OAuthConfig
	oAuthConfig   *oauth2.Config
	endpoint      oauth2.Endpoint
	defaultScopes []string
}

func NewOAuthProvider(name ProviderType, cfg config.OAuthConfig, endpoint oauth2.Endpoint, defaultScopes []string) oAuthProvider {
	return oAuthProvider{
		Name:          name,
		cfg:           cfg,
		endpoint:      endpoint,
		defaultScopes: defaultScopes,
	}
}

func (p *oAuthProvider) GetConfig() *oauth2.Config {
	if p.oAuthConfig == nil {
		p.oAuthConfig = &oauth2.Config{
			ClientID:     p.cfg.Key,
			ClientSecret: p.cfg.Secret,
			Endpoint:     p.endpoint,
			RedirectURL:  fmt.Sprintf("%s/auth/oauth/%s/callback", config.Settings.Service.Fqdn, p.Name),
			Scopes:       append(p.defaultScopes, p.cfg.Scopes...),
		}
	}

	return p.oAuthConfig
}

func (p *oAuthProvider) GetToken(ctx context.Context, code string) (*oauth2.Token, error) {
	cfg := p.GetConfig()

	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange oauth code: %w", err)
	}

	logger.FromContext(ctx).Debug("oauth token received", "token", token)

	return token, nil
}

func (p *oAuthProvider) GetRedirectURL() string {
	cfg := p.GetConfig()

	return cfg.AuthCodeURL(fmt.Sprintf("state:%s", p.Name), oauth2.AccessTypeOnline)
}
