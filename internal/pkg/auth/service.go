package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"recally/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

var ErrUnAuthorized = errors.New("401: username or password or token is invalid")

const DummyUserName = "dummy_user"

// OAuthAdapter interface abstracts OAuth provider implementations
// This interface is implemented by the adapter package to provide OAuth functionality
type OAuthAdapter interface {
	GetAuthURL(ctx context.Context, tx db.DBTX, provider string) (url string, err error)
	HandleCallback(ctx context.Context, tx db.DBTX, provider, code, state string) (oAuth2User any, err error)
	ListProviders() []string
}

type Service struct {
	dao     dto
	adapter OAuthAdapter
}

func New() *Service {
	return &Service{
		dao: db.New(),
	}
}

func NewWithAdapter(adapter OAuthAdapter) *Service {
	return &Service{
		dao:     db.New(),
		adapter: adapter,
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

func (s *Service) GetDummyUser(ctx context.Context, tx db.DBTX) (*UserDTO, error) {
	user, err := s.dao.GetUserByUsername(ctx, tx, pgtype.Text{
		String: DummyUserName,
		Valid:  true,
	})
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

func (s *Service) AuthByPassword(ctx context.Context, tx db.DBTX, email, password string) (*UserDTO, error) {
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

// HandleOAuth2UserLogin processes OAuth login using the Goth adapter
// This method uses the new Goth-based OAuth flow with enhanced security
// including CSRF protection via state validation
//
// Parameters:
//   - ctx: Request context
//   - tx: Database transaction for atomic operations
//   - provider: OAuth provider name (e.g., "google", "github")
//   - code: Authorization code from OAuth callback
//   - state: State parameter for CSRF protection
//
// Returns:
//   - *UserDTO: User information after successful authentication
//   - error: Error if authentication fails at any step
//
// The method performs:
// 1. Validates OAuth callback with state verification (CSRF protection)
// 2. Fetches user info from OAuth provider
// 3. Links to existing user by provider ID or email
// 4. Creates new user if needed
// 5. Updates OAuth connection tokens
func (s *Service) HandleOAuth2UserLogin(ctx context.Context, tx db.DBTX, provider, code, state string) (*UserDTO, error) {
	if s.adapter == nil {
		return nil, fmt.Errorf("OAuth adapter not configured")
	}

	// Use adapter to handle callback with state validation
	oUserInterface, err := s.adapter.HandleCallback(ctx, tx, provider, code, state)
	if err != nil {
		return nil, fmt.Errorf("OAuth callback failed: %w", err)
	}

	// Extract OAuth2User from interface (pattern matching on the returned struct)
	// The adapter returns a struct with fields: Provider, ID, Name, Email, Avatar, AccessToken, RefreshToken, TokenExpiresAt, RawData
	oUserValue := oUserInterface.(struct {
		Provider       string
		ID             string
		Name           string
		Email          string
		Avatar         string
		AccessToken    string
		RefreshToken   string
		TokenExpiresAt time.Time
		RawData        []byte
	})

	user := &UserDTO{}

	// Check if user exists with this OAuth provider
	dbUser, err := s.dao.GetUserByOAuthProviderId(ctx, tx, db.GetUserByOAuthProviderIdParams{
		Provider:       provider,
		ProviderUserID: oUserValue.ID,
	})
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, fmt.Errorf("get user by oauth provider id failed: %w", err)
		}

		// Check if user with the same email exists
		dbUser, err = s.dao.GetUserByEmail(ctx, tx, pgtype.Text{
			String: oUserValue.Email,
			Valid:  oUserValue.Email != "",
		})
		if err != nil {
			if !db.IsNotFoundError(err) {
				return nil, fmt.Errorf("get user by email failed: %w", err)
			}

			// User not found, create new user
			createUserParams := db.CreateUserParams{
				Username: pgtype.Text{String: fmt.Sprintf("%s-%s", oUserValue.Provider, oUserValue.Name), Valid: oUserValue.Name != ""},
				Email:    pgtype.Text{String: oUserValue.Email, Valid: oUserValue.Email != ""},
				Status:   "active",
			}

			dbUser, err = s.dao.CreateUser(ctx, tx, createUserParams)
			if err != nil {
				return nil, fmt.Errorf("create user failed: %w", err)
			}
		}

		// Create OAuth connection
		params := db.CreateOAuthConnectionParams{
			UserID:         dbUser.Uuid,
			Provider:       provider,
			ProviderUserID: oUserValue.ID,
			ProviderEmail:  pgtype.Text{String: oUserValue.Email, Valid: oUserValue.Email != ""},
			AccessToken:    pgtype.Text{String: oUserValue.AccessToken, Valid: oUserValue.AccessToken != ""},
			RefreshToken:   pgtype.Text{String: oUserValue.RefreshToken, Valid: oUserValue.RefreshToken != ""},
			TokenExpiresAt: pgtype.Timestamptz{Time: oUserValue.TokenExpiresAt, Valid: !oUserValue.TokenExpiresAt.IsZero()},
			ProviderData:   oUserValue.RawData,
		}

		_, err = s.dao.CreateOAuthConnection(ctx, tx, params)
		if err != nil {
			return nil, fmt.Errorf("create oauth connection failed: %w", err)
		}
	} else {
		// Update OAuth connection with fresh tokens
		params := db.UpdateOAuthConnectionParams{
			UserID:         dbUser.Uuid,
			Provider:       provider,
			ProviderUserID: oUserValue.ID,
			ProviderEmail:  pgtype.Text{String: oUserValue.Email, Valid: oUserValue.Email != ""},
			AccessToken:    pgtype.Text{String: oUserValue.AccessToken, Valid: oUserValue.AccessToken != ""},
			RefreshToken:   pgtype.Text{String: oUserValue.RefreshToken, Valid: oUserValue.RefreshToken != ""},
			TokenExpiresAt: pgtype.Timestamptz{Time: oUserValue.TokenExpiresAt, Valid: !oUserValue.TokenExpiresAt.IsZero()},
			ProviderData:   oUserValue.RawData,
		}

		_, err = s.dao.UpdateOAuthConnection(ctx, tx, params)
		if err != nil {
			return nil, fmt.Errorf("update oauth connection failed: %w", err)
		}
	}

	user.Load(&dbUser)

	return user, nil
}
