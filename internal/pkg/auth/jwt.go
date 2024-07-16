package auth

import (
	"fmt"
	"time"
	"vibrain/internal/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

var signingMethod = jwt.SigningMethodHS256

type JwtUser struct {
	UserID        string    `json:"user_id"`
	OAuthProvider string    `json:"oauth_provider,omitempty"`
	AccessToken   string    `json:"access_token"`
	TokenType     string    `json:"token_type,omitempty"`
	Expiry        time.Time `json:"expiry,omitempty"`
}

func getJWTSecret() []byte {
	return []byte(config.Settings.JWTSecret)
}

func GenerateJWT(user JwtUser) (string, error) {
	token := jwt.NewWithClaims(signingMethod, jwt.MapClaims{
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"user": user,
	})
	return token.SignedString(getJWTSecret())
}

func ValidateJWT(tokenString string) (*JwtUser, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getJWTSecret(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid jwt token: %w", err)
	}

	claim := token.Claims.(jwt.MapClaims)
	user, ok := claim["user"]
	if !ok {
		return nil, fmt.Errorf("user claim not found")
	}

	return user.(*JwtUser), nil
}
