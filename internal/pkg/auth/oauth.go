package auth

import (
	"context"
	"fmt"
	"recally/internal/pkg/db"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) GetOAuth2RedirectURL(ctx context.Context, provider string) (string, error) {
	oProvider, err := GetOAuthProvider(provider)
	if err != nil {
		return "", err
	}
	return oProvider.GetRedirectURL(), nil
}

// HandleOAuth2Callback processes OAuth2 callback requests after a user authorizes with a provider.
// It handles user authentication and account linking through OAuth providers.
//
// The function performs the following steps:
// 1. Exchanges the authorization code for an access token
// 2. Retrieves user information from the OAuth provider
// 3. Checks if a user already exists with the provider ID
// 4. If no user exists:
//   - Checks for existing user with same email
//   - Creates new user if needed
//   - Creates OAuth connection record
//
// 5. If user exists:
//   - Updates OAuth token information
//
// Parameters:
//   - ctx: Context for the request
//   - tx: Database transaction
//   - provider: OAuth provider name (e.g. "google", "github")
//   - code: Authorization code from OAuth callback
//
// Returns:
//   - *UserDTO: User data transfer object containing user information
//   - error: Error if any step fails
//
// Error cases:
//   - Invalid OAuth provider
//   - Failed to get access token
//   - Failed to get user info
//   - Database errors during user/OAuth operations
func (s *Service) HandleOAuth2Callback(ctx context.Context, tx db.DBTX, provider, code string) (*UserDTO, error) {
	oProvider, err := GetOAuthProvider(provider)
	if err != nil {
		return nil, err
	}
	token, err := oProvider.GetToken(ctx, code)
	if err != nil {
		return nil, err
	}

	oUser, err := oProvider.GetUser(ctx, token)
	if err != nil {
		return nil, err
	}

	user := &UserDTO{}

	dbUser, err := s.dao.GetUserByOAuthProviderId(ctx, tx, db.GetUserByOAuthProviderIdParams{
		Provider:       provider,
		ProviderUserID: oUser.ID,
	})
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, fmt.Errorf("get user by oauth provider id failed: %w", err)
		}

		// check if user with the same email exists
		dbUser, err = s.dao.GetUserByEmail(ctx, tx, pgtype.Text{
			String: oUser.Email,
			Valid:  oUser.Email != "",
		})
		if err != nil {
			if !db.IsNotFoundError(err) {
				return nil, fmt.Errorf("get user by email failed: %w", err)
			}

			// user not found, create new user
			createUserParams := db.CreateUserParams{
				Username: pgtype.Text{String: oUser.Name, Valid: oUser.Name != ""},
				Email:    pgtype.Text{String: oUser.Email, Valid: oUser.Email != ""},
			}
			dbUser, err = s.dao.CreateUser(ctx, tx, createUserParams)
			if err != nil {
				return nil, fmt.Errorf("create user failed: %w", err)
			}
		}

		// create oauth connection
		params := db.CreateOAuthConnectionParams{
			UserID:         dbUser.Uuid,
			Provider:       provider,
			ProviderUserID: oUser.ID,
			ProviderEmail:  pgtype.Text{String: oUser.Email, Valid: oUser.Email != ""},
			AccessToken:    pgtype.Text{String: token.AccessToken, Valid: token.AccessToken != ""},
			RefreshToken:   pgtype.Text{String: token.RefreshToken, Valid: token.RefreshToken != ""},
			TokenExpiresAt: pgtype.Timestamptz{Time: token.Expiry, Valid: true},
		}
		_, err = s.dao.CreateOAuthConnection(ctx, tx, params)
		if err != nil {
			return nil, fmt.Errorf("create oauth connection failed: %w", err)
		}

	} else {
		// update oauth connection
		params := db.UpdateOAuthConnectionParams{
			UserID:         dbUser.Uuid,
			Provider:       provider,
			AccessToken:    pgtype.Text{String: token.AccessToken, Valid: token.AccessToken != ""},
			RefreshToken:   pgtype.Text{String: token.RefreshToken, Valid: token.RefreshToken != ""},
			TokenExpiresAt: pgtype.Timestamptz{Time: token.Expiry, Valid: true},
		}
		_, err = s.dao.UpdateOAuthConnection(ctx, tx, params)
		if err != nil {
			return nil, fmt.Errorf("update oauth connection failed: %w", err)
		}
	}
	user.Load(&dbUser)
	return user, nil
}
