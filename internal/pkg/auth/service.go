package auth

import (
	"context"
	"encoding/json"
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

func (s *Service) UpdateUserSettings(ctx context.Context, tx db.DBTX, settings UserSettings) (*UserDTO, error) {
	return s.UpdateUser(ctx, tx, uuid.Nil, nil, nil, nil, nil, nil, &settings)
}

func (s *Service) UpdateUserSettingsById(ctx context.Context, tx db.DBTX, userId uuid.UUID, settings UserSettings) (*UserDTO, error) {
	return s.UpdateUser(ctx, tx, userId, nil, nil, nil, nil, nil, &settings)
}

func (s *Service) UpdateUserInfo(ctx context.Context, tx db.DBTX, username, email, phone *string) (*UserDTO, error) {
	return s.UpdateUser(ctx, tx, uuid.Nil, username, email, phone, nil, nil, nil)
}

func (s *Service) UpdateUserPassword(ctx context.Context, tx db.DBTX, currentPassword, password string) (*UserDTO, error) {
	user, err := LoadUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.validatePassword(currentPassword, user.Password); err != nil {
		return nil, fmt.Errorf("failed to validate current password: %w", err)
	}

	return s.UpdateUser(ctx, tx, uuid.Nil, nil, nil, nil, &password, nil, nil)
}

func (s *Service) UpdateUserStatus(ctx context.Context, tx db.DBTX, status string) (*UserDTO, error) {
	return s.UpdateUser(ctx, tx, uuid.Nil, nil, nil, nil, nil, &status, nil)
}

func (s *Service) UpdateUserStatusById(ctx context.Context, tx db.DBTX, userId uuid.UUID, status string) (*UserDTO, error) {
	return s.UpdateUser(ctx, tx, userId, nil, nil, nil, nil, &status, nil)
}

func (s *Service) UpdateUser(ctx context.Context, tx db.DBTX, userId uuid.UUID, username, email, phone, password, status *string, settings *UserSettings) (*UserDTO, error) {
	var user *UserDTO
	var err error
	if userId == uuid.Nil {
		user, err = LoadUserFromContext(ctx)
	} else {
		user, err = s.GetUserById(ctx, tx, userId)
	}
	if err != nil {
		return nil, err
	}

	dbUser := user.Dump()
	params := db.UpdateUserByIdParams{
		Uuid:                dbUser.Uuid,
		Username:            dbUser.Username,
		Email:               dbUser.Email,
		Phone:               dbUser.Phone,
		PasswordHash:        dbUser.PasswordHash,
		ActivateAssistantID: dbUser.ActivateAssistantID,
		ActivateThreadID:    dbUser.ActivateThreadID,
		Status:              dbUser.Status,
		Settings:            dbUser.Settings,
	}

	if username != nil {
		params.Username = pgtype.Text{String: *username, Valid: true}
	}
	if email != nil {
		params.Email = pgtype.Text{String: *email, Valid: true}
	}
	if phone != nil {
		params.Phone = pgtype.Text{String: *phone, Valid: true}
	}
	if password != nil {
		hashedPassword, err := s.hashPassword(*password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		params.PasswordHash = pgtype.Text{String: hashedPassword, Valid: true}
	}
	if status != nil {
		params.Status = *status
	}
	if settings != nil {
		newSettings, _ := json.Marshal(settings)
		params.Settings = newSettings
	}

	userModel, err := s.dao.UpdateUserById(ctx, tx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update user settings: %w", err)
	}
	user.Load(&userModel)
	return user, nil
}
