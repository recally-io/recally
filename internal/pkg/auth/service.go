package auth

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Service struct {
	dao dto
}

func New() *Service {
	return &Service{
		dao: db.New(),
	}
}

func getOAuth2Config(provider string) (*oauth2.Config, error) {
	if provider == "github" {
		oauth := config.Settings.OAuths.Github
		return &oauth2.Config{
			ClientID:     oauth.Key,
			ClientSecret: oauth.Secret,
			Endpoint:     github.Endpoint,
			RedirectURL:  fmt.Sprintf("%s/oauth/github/callback", config.Settings.Service.Fqdn),
			Scopes:       []string{"user:email"},
		}, nil
	}

	return nil, fmt.Errorf("oauth provider '%s' not found", provider)
}

func (s *Service) GetOAuth2RedirectURL(ctx context.Context, provider string) (string, error) {
	cfg, err := getOAuth2Config(provider)
	if err != nil {
		return "", fmt.Errorf("failed to get oauth config: %w", err)
	}
	authCodeUrl := cfg.AuthCodeURL("state:"+provider, oauth2.AccessTypeOnline)
	return authCodeUrl, nil
}

func (s *Service) GetOAuth2Token(ctx context.Context, provider, code string) (*oauth2.Token, error) {
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

func (s *Service) GetUserById(ctx context.Context, tx db.DBTX, userId uuid.UUID) (*UserDTO, error) {
	user, err := s.dao.GetUserById(ctx, tx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	u := new(UserDTO)
	u.Load(&user)
	return u, nil
}

func (s *Service) CreateUser(ctx context.Context, tx db.DBTX, user *UserDTO) (*UserDTO, error) {
	dbUser := user.Dump()
	params := db.CreateUserParams{
		Username:            dbUser.Username,
		Email:               dbUser.Email,
		Github:              dbUser.Github,
		Google:              dbUser.Google,
		Telegram:            dbUser.Telegram,
		ActivateAssistantID: dbUser.ActivateAssistantID,
		ActivateThreadID:    dbUser.ActivateThreadID,
		Status:              dbUser.Status,
	}
	userModel, err := s.dao.CreateUser(ctx, tx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	user.Load(&userModel)
	return user, nil
}
