package auth

import (
	"context"
	"errors"
	"fmt"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

var ErrUnAuthorized = errors.New("401: username or password or token is invalid")

type Service struct {
	dao dto
}

func New() *Service {
	return &Service{
		dao: db.New(),
	}
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
	if user.Password != "" {
		hashedPassword, err := s.hashPassword(user.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}
	dbUser := user.Dump()
	params := db.CreateUserParams{
		Username:            dbUser.Username,
		Email:               dbUser.Email,
		PasswordHash:        dbUser.PasswordHash,
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

func (s *Service) AuthByPassword(ctx context.Context, tx db.DBTX, email string, password string) (*UserDTO, error) {
	user, err := s.dao.GetUserByEmail(ctx, tx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		if db.IsNotFoundError(err) {
			return nil, ErrUnAuthorized
		}
		return nil, fmt.Errorf("failed to load user: %w", err)
	}
	if err := s.validatePassword(password, user.PasswordHash.String); err != nil {
		return nil, ErrUnAuthorized
	}

	u := new(UserDTO)
	u.Load(&user)
	return u, nil
}

func (s *Service) GetTelegramUser(ctx context.Context, tx db.DBTX, userID string) (*UserDTO, error) {
	user, err := s.dao.GetUserByOAuthProviderId(ctx, tx, db.GetUserByOAuthProviderIdParams{
		Provider:       "telegram",
		ProviderUserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get telegram user: %w", err)
	}
	u := new(UserDTO)
	u.Load(&user)
	return u, nil
}

func (s *Service) CreateTelegramUser(ctx context.Context, tx db.DBTX, userName, userID string) (*UserDTO, error) {
	// create user
	user, err := s.dao.CreateUser(ctx, tx, db.CreateUserParams{
		Username: pgtype.Text{String: fmt.Sprintf("TG-%s", userName), Valid: true},
		Status:   "active",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram user: %w", err)
	}
	// create oauth telegram connection
	_, err = s.dao.CreateOAuthConnection(ctx, tx, db.CreateOAuthConnectionParams{
		UserID:         user.Uuid,
		Provider:       "telegram",
		ProviderUserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram oauth connection: %w", err)
	}
	u := new(UserDTO)
	u.Load(&user)
	return u, nil
}

func (s *Service) UpdateTelegramUser(ctx context.Context, tx db.DBTX, user *UserDTO) (*UserDTO, error) {
	dbUser := user.Dump()
	params := db.UpdateUserByIdParams{
		ActivateAssistantID: dbUser.ActivateAssistantID,
		ActivateThreadID:    dbUser.ActivateThreadID,
	}
	userModel, err := s.dao.UpdateUserById(ctx, tx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update telegram user: %w", err)
	}
	user.Load(&userModel)
	return user, nil
}

func (s *Service) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("auth: failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func (s *Service) validatePassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return ErrUnAuthorized
	}
	return nil
}
