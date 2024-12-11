package auth

import (
	"context"
	"time"
)

// User represents the authenticated user information
type User struct {
	ID                  string
	Username            *string
	Email               *string
	Phone               *string
	FullName            *string
	AvatarURL           *string
	IsActive            bool
	IsVerified          bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
	LastLoginAt         *time.Time
	FailedLoginAttempts int16
	LockedUntil         *time.Time
	Settings            map[string]interface{}
}

// OAuthConnection represents a user's OAuth provider connection
type OAuthConnection struct {
	Provider       string
	ProviderUserID string
	ProviderEmail  *string
	AccessToken    string
	RefreshToken   *string
	TokenExpiresAt *time.Time
	ProviderData   map[string]interface{}
}

// VerificationType represents different types of verification
type VerificationType string

const (
	VerificationTypeEmail  VerificationType = "email"
	VerificationTypePhone  VerificationType = "phone"
	VerificationTypeReset  VerificationType = "reset"
	VerificationTypeTwoFac VerificationType = "2fa"
)

type AuthService interface {
	// User Registration and Authentication
	SignUp(ctx context.Context, username, email, phone, password, fullName string) (*User, error)
	SignIn(ctx context.Context, identifier, password string) (*User, string, error) // identifier can be email/phone/username
	SignOut(ctx context.Context, userID string) error

	// OAuth Operations
	SignUpWithOAuth(ctx context.Context, provider string, code string) (*User, string, error)
	SignInWithOAuth(ctx context.Context, provider string, code string) (*User, string, error)
	LinkOAuthProvider(ctx context.Context, userID string, provider string, code string) error
	UnlinkOAuthProvider(ctx context.Context, userID string, provider string) error
	GetOAuthConnections(ctx context.Context, userID string) ([]OAuthConnection, error)

	// JWT Token Operations
	GenerateToken(ctx context.Context, userID string) (string, error)
	ValidateToken(ctx context.Context, token string) (*User, error)
	RevokeToken(ctx context.Context, token string) error

	// API Key Operations
	GenerateAPIKey(ctx context.Context, userID string, description string) (string, error)
	ValidateAPIKey(ctx context.Context, apiKey string) (*User, error)
	RevokeAPIKey(ctx context.Context, apiKey string) error
	ListAPIKeys(ctx context.Context, userID string) ([]string, error)

	// Verification Operations
	SendVerification(ctx context.Context, userID string, verificationType VerificationType) error
	VerifyToken(ctx context.Context, userID string, token string, verificationType VerificationType) error

	// User Management
	GetUser(ctx context.Context, userID string) (*User, error)
	UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error
	DeactivateUser(ctx context.Context, userID string) error
	ReactivateUser(ctx context.Context, userID string) error

	// Password Management
	ChangePassword(ctx context.Context, userID string, oldPassword, newPassword string) error
	ResetPassword(ctx context.Context, token string, newPassword string) error
	RequestPasswordReset(ctx context.Context, email string) error
}
