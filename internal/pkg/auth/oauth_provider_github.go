package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"recally/internal/pkg/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const GithubProvider ProviderType = "github"

type OAuthProviderGithub struct {
	oAuthProvider
}

// GitHub user structure.
type GitHubUser struct {
	Login       string `json:"login"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar_url"`
	Location    string `json:"location"`
	Bio         string `json:"bio"`
	PublicRepos int    `json:"public_repos"`
}

func NewOAuthProviderGithub(cfg config.OAuthConfig) *OAuthProviderGithub {
	return &OAuthProviderGithub{
		oAuthProvider: NewOAuthProvider(GithubProvider, cfg, github.Endpoint, []string{"user"}),
	}
}

func (p *OAuthProviderGithub) GetUser(ctx context.Context, token *oauth2.Token) (OAuth2User, error) {
	client := p.GetConfig().Client(ctx, token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return OAuth2User{}, fmt.Errorf("failed to get user from github: %w", err)
	}

	defer resp.Body.Close()

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return OAuth2User{}, fmt.Errorf("failed to decode github user: %w", err)
	}

	return OAuth2User{
		Provider:       "github",
		ID:             user.Login,
		Name:           user.Name,
		Email:          user.Email,
		Avatar:         user.Avatar,
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		TokenExpiresAt: token.Expiry,
		RawData:        nil,
	}, nil
}
