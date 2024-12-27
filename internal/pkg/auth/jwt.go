package auth

import (
	"context"
	"fmt"
	"time"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/db"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSigningMethod = jwt.SigningMethodHS256

type JwtUser struct {
	ID          string    `json:"id"`
	AccessToken string    `json:"access_token"`
	Expiry      time.Time `json:"expiry"`
}

func getJWTSecret() []byte {
	return []byte(config.Settings.JWTSecret)
}

func (s *Service) GenerateJWT(userId uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwtSigningMethod, jwt.MapClaims{
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"user_id": userId,
	})
	return token.SignedString(getJWTSecret())
}

func (s *Service) ValidateJWT(ctx context.Context, tx db.DBTX, tokenString string) (*UserDTO, int64, error) {
	userId, exp, err := ValidateJWT(tokenString)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid token: %w", err)
	}
	user, err := s.GetUserById(ctx, tx, userId)
	return user, exp, err
}

func ValidateJWT(tokenString string) (uuid.UUID, int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method: %v", token.Header["alg"])
		}
		return getJWTSecret(), nil
	})
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("jwt: invalid jwt token: %w", err)
	}

	claim := token.Claims.(jwt.MapClaims)
	userId, ok := claim["user_id"]
	if !ok {
		return uuid.Nil, 0, fmt.Errorf("jwt: user claim not found")
	}
	exp := int64(claim["exp"].(float64))

	return uuid.MustParse(userId.(string)), exp, nil
}
