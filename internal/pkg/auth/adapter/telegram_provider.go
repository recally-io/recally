package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/markbates/goth"

	"recally/internal/pkg/config"
)

// TelegramProvider implements goth.Provider for Telegram
type TelegramProvider struct {
	botToken string
	callback string
}

func NewTelegramProvider(cfg config.TelegramConfig, callbackURL string) *TelegramProvider {
	return &TelegramProvider{
		botToken: cfg.Token,
		callback: callbackURL,
	}
}

func (p *TelegramProvider) Name() string {
	return "telegram"
}

func (p *TelegramProvider) SetName(name string) {
	// No-op: Telegram provider name is fixed
}

func (p *TelegramProvider) Debug(debug bool) {
	// No-op: Debug mode not implemented for Telegram provider
}

func (p *TelegramProvider) BeginAuth(state string) (goth.Session, error) {
	return &TelegramSession{
		AuthURL: fmt.Sprintf("https://oauth.telegram.org/auth?bot_id=%s&origin=%s&state=%s",
			p.botToken, p.callback, state),
		State: state,
	}, nil
}

func (p *TelegramProvider) FetchUser(session goth.Session) (goth.User, error) {
	// sess := session.(*TelegramSession)
	// Note: In production, you would use sess to get user-specific auth data

	// Verify Telegram auth data
	// Note: In a real implementation, you would validate the auth data hash
	// using the bot token to ensure the data came from Telegram
	resp, err := http.Get(fmt.Sprintf(
		"https://api.telegram.org/bot%s/getMe",
		p.botToken,
	))
	if err != nil {
		return goth.User{}, fmt.Errorf("failed to verify bot: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var result struct {
		OK     bool `json:"ok"`
		Result struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			PhotoURL  string `json:"photo_url"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return goth.User{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.OK {
		return goth.User{}, errors.New("telegram auth failed")
	}

	// In a real implementation, you would get the actual user data
	// from the callback parameters, not from getMe (which returns bot info)
	// This is a simplified version for demonstration
	return goth.User{
		Provider:  "telegram",
		UserID:    fmt.Sprintf("%d", result.Result.ID),
		Name:      result.Result.FirstName,
		NickName:  result.Result.Username,
		AvatarURL: result.Result.PhotoURL,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Telegram tokens typically don't expire
	}, nil
}

func (p *TelegramProvider) UnmarshalSession(data string) (goth.Session, error) {
	sess := &TelegramSession{}
	err := json.Unmarshal([]byte(data), sess)
	return sess, err
}

// TelegramSession represents a Telegram OAuth session
type TelegramSession struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
	Code    string `json:"code,omitempty"`
}

func (s *TelegramSession) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New("auth URL not set")
	}
	return s.AuthURL, nil
}

func (s *TelegramSession) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	code := params.Get("code")
	if code == "" {
		return "", errors.New("authorization code not found")
	}
	s.Code = code
	return code, nil
}

func (s *TelegramSession) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}
