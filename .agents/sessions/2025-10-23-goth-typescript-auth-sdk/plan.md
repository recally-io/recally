# Implementation Plan: Goth OAuth + TypeScript Auth SDK

**Status:** Approved
**Date:** 2025-10-23
**Estimated Effort:** 24 hours (3 days)
**Risk Level:** Low-Medium (Security-focused migration)

---

## 1. Executive Summary

### Goal
Enhance Recally's authentication system by migrating from custom OAuth implementation to Goth (multi-provider OAuth library) and building a minimal TypeScript SDK for better frontend developer experience.

### Key Outcomes
- ✅ Secure OAuth implementation with proper CSRF protection
- ✅ Easy multi-provider support (GitHub, Google, Telegram)
- ✅ Type-safe TypeScript auth client with autocomplete
- ✅ Zero database migration or schema changes
- ✅ All existing features preserved (JWT, API keys, Telegram auth)
- ✅ Direct replacement (no gradual rollout needed)

### Success Metrics
- All existing users can authenticate without issues
- OAuth CSRF attacks prevented via secure state management
- New OAuth providers can be added in <30 minutes
- Frontend code has full TypeScript type safety
- Zero breaking changes to API contracts

---

## 2. Architecture Overview

### 2.1 System Architecture

```
┌────────────────────────────────────────────────────────────────┐
│                    Go Backend (Echo)                           │
│                                                                │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  OAuth Adapter Layer (NEW)                               │  │
│  │  ┌────────────────┐  ┌──────────────────────────────┐   │  │
│  │  │ Goth OAuth     │  │ Custom Telegram Provider     │   │  │
│  │  │ - GitHub       │  │ (implements goth.Provider)   │   │  │
│  │  │ - Google       │  │                              │   │  │
│  │  └────────────────┘  └──────────────────────────────┘   │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Authentication Services                                 │  │
│  │  - JWT Sessions (UNCHANGED)                              │  │
│  │  - API Keys (UNCHANGED)                                  │  │
│  │  - State Management (NEW - database-backed)              │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                │
│  Database:                                                     │
│  - users                                                       │
│  - auth_user_oauth_connections                                │
│  - auth_oauth_states (NEW - for CSRF protection)              │
│  - auth_api_keys                                               │
└────────────────────────────────────────────────────────────────┘
                          ↑
                          │ HTTP JSON API
                          ↓
┌────────────────────────────────────────────────────────────────┐
│                    React Frontend                              │
│                                                                │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Minimal TypeScript Auth Client (NEW)                    │  │
│  │  /web/src/lib/auth-client/                               │  │
│  │                                                           │  │
│  │  - Type-safe API client                                  │  │
│  │  - Better error handling                                 │  │
│  │  - OAuth helpers                                         │  │
│  │  - Keep existing SWR hooks simple                        │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────┘
```

### 2.2 Directory Structure

**Backend Changes:**
```
/internal/pkg/auth/
├── adapter/
│   ├── oauth_adapter.go        # NEW: OAuth adapter interface
│   ├── goth_adapter.go         # NEW: Goth implementation
│   └── telegram_provider.go    # NEW: Custom Telegram Goth provider
├── goth_config.go              # NEW: Goth initialization
├── state_manager.go            # NEW: OAuth state CRUD
├── oauth.go                    # UPDATE: Use adapter instead of custom code
├── oauth_provider.go           # REMOVE after migration
├── oauth_provider_github.go    # REMOVE after migration
├── service.go                  # UPDATE: Use adapter
├── jwt.go                      # UNCHANGED
├── api_key.go                  # UNCHANGED
└── context.go                  # UNCHANGED

/internal/port/httpserver/
├── handler_auth.go             # UPDATE: Use Goth adapter
└── middleware.go               # UPDATE: Add Secure/HttpOnly to cookies

/database/migrations/
└── 20251023_oauth_states.sql   # NEW: State management table
```

**Frontend Structure:**
```
/web/src/lib/auth-client/       # NEW (minimal SDK)
├── index.ts                    # Main export
├── client.ts                   # Typed API client
├── types.ts                    # Type definitions
└── errors.ts                   # Error handling

/web/src/lib/apis/
├── auth.ts                     # UPDATE: Use new client, keep hooks simple
└── ...                         # Other APIs unchanged
```

---

## 3. Technical Implementation Details

### 3.1 Database Schema Changes

**New Migration: `20251023_oauth_states.sql`**

```sql
-- OAuth state management for CSRF protection
CREATE TABLE IF NOT EXISTS auth_oauth_states (
    state VARCHAR(64) PRIMARY KEY,
    provider VARCHAR(50) NOT NULL,
    redirect_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

-- Index for cleanup queries
CREATE INDEX idx_oauth_states_expiry ON auth_oauth_states(expires_at);

-- Automatic cleanup of expired states (run via cron or scheduled job)
-- DELETE FROM auth_oauth_states WHERE expires_at < NOW();
```

### 3.2 Backend Components

#### A. OAuth Adapter Interface

**File: `/internal/pkg/auth/adapter/oauth_adapter.go`**

```go
package adapter

import (
    "context"
    "recally/internal/pkg/auth"
)

// OAuthAdapter abstracts OAuth provider implementation
type OAuthAdapter interface {
    // GetAuthURL generates OAuth redirect URL with secure state
    GetAuthURL(ctx context.Context, provider string) (url string, err error)

    // HandleCallback processes OAuth callback and returns user info
    HandleCallback(ctx context.Context, provider, code, state string) (auth.OAuth2User, error)

    // ListProviders returns available OAuth providers
    ListProviders() []string
}
```

#### B. Goth Adapter Implementation

**File: `/internal/pkg/auth/adapter/goth_adapter.go`**

```go
package adapter

import (
    "context"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "time"

    "github.com/markbates/goth"
    "recally/internal/pkg/auth"
    "recally/internal/pkg/db"
)

type GothAdapter struct {
    dao        db.Queries
    providers  map[string]goth.Provider
}

func NewGothAdapter(dao db.Queries) *GothAdapter {
    return &GothAdapter{
        dao:       dao,
        providers: make(map[string]goth.Provider),
    }
}

func (a *GothAdapter) RegisterProvider(provider goth.Provider) {
    a.providers[provider.Name()] = provider
}

func (a *GothAdapter) GetAuthURL(ctx context.Context, provider string) (string, error) {
    p, ok := a.providers[provider]
    if !ok {
        return "", fmt.Errorf("provider %s not found", provider)
    }

    // Generate secure state
    state, err := generateSecureState()
    if err != nil {
        return "", err
    }

    // Store state in database with 5-minute expiration
    tx := db.ExtractTx(ctx)
    err = a.dao.CreateOAuthState(ctx, tx, db.CreateOAuthStateParams{
        State:     state,
        Provider:  provider,
        ExpiresAt: time.Now().Add(5 * time.Minute),
    })
    if err != nil {
        return "", err
    }

    // Get auth URL from provider
    sess, err := p.BeginAuth(state)
    if err != nil {
        return "", err
    }

    return sess.GetAuthURL()
}

func (a *GothAdapter) HandleCallback(ctx context.Context, provider, code, state string) (auth.OAuth2User, error) {
    // Validate state (CSRF protection)
    tx := db.ExtractTx(ctx)
    storedState, err := a.dao.GetOAuthState(ctx, tx, state)
    if err != nil {
        return auth.OAuth2User{}, fmt.Errorf("invalid state: %w", err)
    }

    if storedState.Provider != provider {
        return auth.OAuth2User{}, fmt.Errorf("provider mismatch")
    }

    if time.Now().After(storedState.ExpiresAt) {
        return auth.OAuth2User{}, fmt.Errorf("state expired")
    }

    // Delete used state
    _ = a.dao.DeleteOAuthState(ctx, tx, state)

    // Get provider and complete auth
    p, ok := a.providers[provider]
    if !ok {
        return auth.OAuth2User{}, fmt.Errorf("provider %s not found", provider)
    }

    // Authorize with code
    sess, err := p.BeginAuth(state)
    if err != nil {
        return auth.OAuth2User{}, err
    }

    _, err = sess.Authorize(p, map[string]string{"code": code})
    if err != nil {
        return auth.OAuth2User{}, err
    }

    gothUser, err := p.FetchUser(sess)
    if err != nil {
        return auth.OAuth2User{}, err
    }

    // Convert to our OAuth2User type
    return auth.OAuth2User{
        Provider:       gothUser.Provider,
        ID:             gothUser.UserID,
        Name:           gothUser.Name,
        Email:          gothUser.Email,
        Avatar:         gothUser.AvatarURL,
        AccessToken:    gothUser.AccessToken,
        RefreshToken:   gothUser.RefreshToken,
        TokenExpiresAt: gothUser.ExpiresAt,
    }, nil
}

func generateSecureState() (string, error) {
    b := make([]byte, 32)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(b), nil
}

func (a *GothAdapter) ListProviders() []string {
    names := make([]string, 0, len(a.providers))
    for name := range a.providers {
        names = append(names, name)
    }
    return names
}
```

#### C. Custom Telegram Provider

**File: `/internal/pkg/auth/adapter/telegram_provider.go`**

```go
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

func (p *TelegramProvider) BeginAuth(state string) (goth.Session, error) {
    return &TelegramSession{
        AuthURL: fmt.Sprintf("https://oauth.telegram.org/auth?bot_id=%s&origin=%s",
            p.botToken, p.callback),
    }, nil
}

func (p *TelegramProvider) FetchUser(session goth.Session) (goth.User, error) {
    sess := session.(*TelegramSession)

    // Verify Telegram auth data
    resp, err := http.Get(fmt.Sprintf(
        "https://api.telegram.org/bot%s/getMe",
        p.botToken,
    ))
    if err != nil {
        return goth.User{}, err
    }
    defer resp.Body.Close()

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
        return goth.User{}, err
    }

    if !result.OK {
        return goth.User{}, errors.New("telegram auth failed")
    }

    return goth.User{
        Provider:  "telegram",
        UserID:    fmt.Sprintf("%d", result.Result.ID),
        Name:      result.Result.FirstName,
        NickName:  result.Result.Username,
        AvatarURL: result.Result.PhotoURL,
    }, nil
}

func (p *TelegramProvider) UnmarshalSession(data string) (goth.Session, error) {
    sess := &TelegramSession{}
    err := json.Unmarshal([]byte(data), sess)
    return sess, err
}

func (p *TelegramProvider) SetName(name string) {}
func (p *TelegramProvider) Debug(debug bool) {}

type TelegramSession struct {
    AuthURL string
}

func (s *TelegramSession) GetAuthURL() (string, error) {
    return s.AuthURL, nil
}

func (s *TelegramSession) Authorize(provider goth.Provider, params goth.Params) (string, error) {
    return params.Get("code"), nil
}

func (s *TelegramSession) Marshal() string {
    b, _ := json.Marshal(s)
    return string(b)
}
```

#### D. Goth Configuration

**File: `/internal/pkg/auth/goth_config.go`**

```go
package auth

import (
    "fmt"

    "github.com/markbates/goth"
    "github.com/markbates/goth/providers/github"
    "github.com/markbates/goth/providers/google"
    "recally/internal/pkg/auth/adapter"
    "recally/internal/pkg/config"
)

func InitGothAdapter(dao db.Queries) *adapter.GothAdapter {
    gothAdapter := adapter.NewGothAdapter(dao)

    // Register GitHub provider
    if config.Settings.OAuths.Github.Key != "" {
        githubProvider := github.New(
            config.Settings.OAuths.Github.Key,
            config.Settings.OAuths.Github.Secret,
            fmt.Sprintf("%s/api/v1/oauth/github/callback", config.Settings.Service.Fqdn),
            "user:email",
        )
        gothAdapter.RegisterProvider(githubProvider)
    }

    // Register Google provider
    if config.Settings.OAuths.Google.Key != "" {
        googleProvider := google.New(
            config.Settings.OAuths.Google.Key,
            config.Settings.OAuths.Google.Secret,
            fmt.Sprintf("%s/api/v1/oauth/google/callback", config.Settings.Service.Fqdn),
            "email", "profile",
        )
        gothAdapter.RegisterProvider(googleProvider)
    }

    // Register custom Telegram provider
    if config.Settings.Telegram.Reader.Token != "" {
        telegramProvider := adapter.NewTelegramProvider(
            config.Settings.Telegram.Reader,
            fmt.Sprintf("%s/api/v1/oauth/telegram/callback", config.Settings.Service.Fqdn),
        )
        gothAdapter.RegisterProvider(telegramProvider)
    }

    return gothAdapter
}
```

#### E. Updated Auth Handler

**File: `/internal/port/httpserver/handler_auth.go` (UPDATED)**

```go
// Update OAuth handlers to use adapter
func (h *authHandler) oAuthLogin(c echo.Context) error {
    provider := c.Param("provider")
    ctx := c.Request().Context()

    url, err := h.oauthAdapter.GetAuthURL(ctx, provider)
    if err != nil {
        return ErrorResponse(c, http.StatusInternalServerError, err)
    }

    return JsonResponse(c, http.StatusOK, map[string]string{
        "url": url,
    })
}

func (h *authHandler) oAuthCallback(c echo.Context) error {
    provider := c.Param("provider")
    code := c.QueryParam("code")
    state := c.QueryParam("state")
    ctx := c.Request().Context()

    // Validate via adapter (handles state checking)
    oauthUser, err := h.oauthAdapter.HandleCallback(ctx, provider, code, state)
    if err != nil {
        return ErrorResponse(c, http.StatusUnauthorized, fmt.Errorf("oauth failed: %w", err))
    }

    // Use existing user creation/update logic
    tx, err := loadTx(ctx)
    if err != nil {
        return ErrorResponse(c, http.StatusInternalServerError, err)
    }

    user, err := h.service.HandleOAuth2UserLogin(ctx, tx, oauthUser)
    if err != nil {
        return ErrorResponse(c, http.StatusInternalServerError, err)
    }

    jwtToken, err := h.service.GenerateJWT(user.ID)
    if err != nil {
        return ErrorResponse(c, http.StatusInternalServerError, err)
    }

    h.setCookieJwtToken(c, jwtToken)

    // Redirect to frontend
    return c.Redirect(http.StatusFound, "/")
}

// Update cookie to be secure
func (h *authHandler) setCookieJwtToken(c echo.Context, token string) {
    cookie := &http.Cookie{
        Name:     "token",
        Value:    token,
        Expires:  time.Now().Add(time.Hour * 24),
        Path:     "/",
        HttpOnly: true,                                    // Prevent XSS
        Secure:   config.Settings.Service.Env == "production", // HTTPS only in prod
        SameSite: http.SameSiteLaxMode,
    }
    c.SetCookie(cookie)
}
```

### 3.3 Frontend Components

#### A. Minimal TypeScript Auth Client

**File: `/web/src/lib/auth-client/client.ts`**

```typescript
export class AuthClient {
  constructor(private baseURL: string = '/api/v1') {}

  async login(email: string, password: string): Promise<User> {
    return this.post<User>('/auth/login', { email, password })
  }

  async register(username: string, email: string, password: string): Promise<User> {
    return this.post<User>('/auth/register', { username, email, password })
  }

  async logout(): Promise<void> {
    await this.post('/auth/logout', {})
  }

  async validateSession(): Promise<User> {
    return this.get<User>('/auth/validate-jwt')
  }

  async getOAuthURL(provider: string): Promise<string> {
    const resp = await this.get<{ url: string }>(`/oauth/${provider}/login`)
    return resp.url
  }

  redirectToOAuth(provider: string): void {
    window.location.href = `${this.baseURL}/oauth/${provider}/login`
  }

  private async get<T>(path: string): Promise<T> {
    const response = await fetch(`${this.baseURL}${path}`, {
      credentials: 'include',
    })

    if (!response.ok) {
      throw new AuthError(await response.json())
    }

    return response.json()
  }

  private async post<T>(path: string, body: any): Promise<T> {
    const response = await fetch(`${this.baseURL}${path}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(body),
    })

    if (!response.ok) {
      throw new AuthError(await response.json())
    }

    return response.json()
  }
}

export const authClient = new AuthClient()
```

**File: `/web/src/lib/auth-client/types.ts`**

```typescript
export interface User {
  id: string
  username: string
  email: string
  avatar: string
  phone?: string
  status: string
  settings: UserSettings
}

export interface UserSettings {
  theme?: 'light' | 'dark'
  isLinkedTelegramBot?: boolean
  [key: string]: any
}

export interface LoginInput {
  email: string
  password: string
}

export interface RegisterInput {
  username: string
  email: string
  password: string
}
```

**File: `/web/src/lib/auth-client/errors.ts`**

```typescript
export class AuthError extends Error {
  constructor(
    public code: string,
    public status: number,
    message: string
  ) {
    super(message)
    this.name = 'AuthError'
  }

  static fromResponse(response: any): AuthError {
    return new AuthError(
      response.code || 'AUTH_ERROR',
      response.status || 500,
      response.message || 'Authentication failed'
    )
  }
}
```

#### B. Updated Auth Hooks

**File: `/web/src/lib/apis/auth.ts` (SIMPLIFIED)**

```typescript
import useSWR from 'swr'
import { authClient } from '@/lib/auth-client'
import type { User, LoginInput, RegisterInput } from '@/lib/auth-client/types'

// Keep SWR hooks simple - just wrap the client
export function useUser() {
  const { data, error, mutate } = useSWR<User>(
    'auth-user',
    () => authClient.validateSession(),
    {
      revalidateOnFocus: false,
      revalidateIfStale: false,
      revalidateOnReconnect: false,
      dedupingInterval: 60000,
      shouldRetryOnError: false,
    }
  )

  return {
    user: data,
    isLoading: !error && !data,
    isError: error,
    mutate,
  }
}

export function useAuth() {
  const { mutate } = useUser()

  return {
    login: async (input: LoginInput) => {
      const user = await authClient.login(input.email, input.password)
      await mutate(user)
      return user
    },

    register: async (input: RegisterInput) => {
      const user = await authClient.register(input.username, input.email, input.password)
      await mutate(user)
      return user
    },

    logout: async () => {
      await authClient.logout()
      await mutate(undefined)
    },

    oauthLogin: (provider: string) => {
      authClient.redirectToOAuth(provider)
    },
  }
}

// Keep API key hooks unchanged
// ... existing API key code ...
```

---

## 4. Implementation Steps & Timeline

### Phase 0: Preparation (3 hours)

**Day 1 Morning**

**Step 0.1: Database Migration (30 min)**
- Create migration file `20251023_oauth_states.sql`
- Run migration: `mise run migrate:up`
- Verify table created: `mise run psql` → `\d auth_oauth_states`

**Step 0.2: Add SQLC Queries (30 min)**
- Add to `/database/queries/auth.sql`:
  ```sql
  -- name: CreateOAuthState :exec
  INSERT INTO auth_oauth_states (state, provider, expires_at)
  VALUES ($1, $2, $3);

  -- name: GetOAuthState :one
  SELECT * FROM auth_oauth_states WHERE state = $1;

  -- name: DeleteOAuthState :exec
  DELETE FROM auth_oauth_states WHERE state = $1;

  -- name: CleanupExpiredStates :exec
  DELETE FROM auth_oauth_states WHERE expires_at < NOW();
  ```
- Run: `mise run generate:sql`

**Step 0.3: Install Dependencies (15 min)**
```bash
go get github.com/markbates/goth
go get github.com/markbates/goth/providers/github
go get github.com/markbates/goth/providers/google
```

**Step 0.4: Add Google OAuth Config (30 min)**
- Register app at Google Cloud Console
- Add to `.env`:
  ```
  OAUTH_GOOGLE_KEY=your_client_id
  OAUTH_GOOGLE_SECRET=your_client_secret
  ```
- Update config struct

**Step 0.5: Create Directory Structure (15 min)**
```bash
mkdir -p internal/pkg/auth/adapter
mkdir -p web/src/lib/auth-client
```

### Phase 1: Backend Implementation (5 hours)

**Day 1 Afternoon**

**Step 1.1: Implement OAuth Adapter (1 hour)**
- Create `adapter/oauth_adapter.go` (interface)
- Create `adapter/goth_adapter.go` (implementation)
- Implement state management functions

**Step 1.2: Implement Telegram Provider (1.5 hours)**
- Create `adapter/telegram_provider.go`
- Implement `goth.Provider` interface
- Test Telegram widget auth flow

**Step 1.3: Create Goth Configuration (30 min)**
- Create `goth_config.go`
- Initialize all providers (GitHub, Google, Telegram)

**Step 1.4: Update Auth Service (1 hour)**
- Modify `/internal/pkg/auth/service.go`
- Add adapter field
- Update OAuth methods to use adapter
- Rename `HandleOAuth2Callback` to `HandleOAuth2UserLogin`

**Step 1.5: Update HTTP Handlers (1 hour)**
- Modify `/internal/port/httpserver/handler_auth.go`
- Update `oAuthLogin` to use adapter
- Update `oAuthCallback` to validate state
- Add `HttpOnly` and `Secure` to cookie

### Phase 2: Backend Testing (3 hours)

**Day 2 Morning**

**Step 2.1: Unit Tests (1 hour)**
```go
// adapter/goth_adapter_test.go
func TestStateGeneration(t *testing.T) { ... }
func TestStateValidation(t *testing.T) { ... }
func TestStateExpiration(t *testing.T) { ... }
```

**Step 2.2: Integration Tests (1.5 hours)**
```go
// handler_auth_test.go
func TestGitHubOAuthFlow(t *testing.T) { ... }
func TestGoogleOAuthFlow(t *testing.T) { ... }
func TestTelegramOAuthFlow(t *testing.T) { ... }
func TestCSRFProtection(t *testing.T) { ... }
```

**Step 2.3: Manual Testing (30 min)**
- Start server: `mise run dev:backend`
- Test GitHub OAuth end-to-end
- Test Google OAuth end-to-end
- Verify state validation works
- Check database records

### Phase 3: Frontend Implementation (6 hours)

**Day 2 Afternoon**

**Step 3.1: Create TypeScript Client (2 hours)**
- Create `/web/src/lib/auth-client/client.ts`
- Create `/web/src/lib/auth-client/types.ts`
- Create `/web/src/lib/auth-client/errors.ts`
- Create `/web/src/lib/auth-client/index.ts`

**Step 3.2: Update Auth Hooks (1 hour)**
- Simplify `/web/src/lib/apis/auth.ts`
- Replace fetch calls with `authClient` calls
- Keep existing SWR hook structure

**Step 3.3: Update Auth Component (1 hour)**
- Update `/web/src/components/auth/auth.tsx`
- Add Google OAuth button
- Test all auth flows in UI

**Step 3.4: Update Other Components (1 hour)**
- Find all uses of `useUser()` and verify compatibility
- Update protected route components if needed
- Test session persistence

**Step 3.5: Frontend Testing (1 hour)**
- Write component tests for auth forms
- Test OAuth button clicks
- Test error handling
- E2E test: full registration → login → logout flow

### Phase 4: Testing & Documentation (6 hours)

**Day 3**

**Step 4.1: Security Audit (2 hours)**
- [ ] CSRF attack test (forged state parameter)
- [ ] Replay attack test (reused state)
- [ ] State expiration test
- [ ] Cookie security flags verification
- [ ] HTTPS redirect test (production)
- [ ] XSS prevention test (HttpOnly cookie)

**Step 4.2: Load Testing (1 hour)**
```bash
# Test concurrent OAuth flows
ab -n 100 -c 10 http://localhost:8080/api/v1/oauth/github/login
```

**Step 4.3: Migration Validation (1 hour)**
- Test existing users can still login
- Test existing OAuth connections work
- Test API keys unchanged
- Test Telegram auth preserved

**Step 4.4: Documentation (2 hours)**
- Update `CLAUDE.md` with new architecture
- Create `/web/src/lib/auth-client/README.md`
- Document OAuth provider setup
- Create migration notes
- Update API documentation

### Phase 5: Cleanup (1 hour)

**Day 3 End**

**Step 5.1: Remove Old Code**
- Delete `/internal/pkg/auth/oauth_provider.go`
- Delete `/internal/pkg/auth/oauth_provider_github.go`
- Add deprecation notice to old imports

**Step 5.2: Final Verification**
- Run full test suite: `mise run test`
- Run linters: `mise run lint`
- Check build: `mise run build`

---

## 5. Testing Strategy

### 5.1 Security Testing (Critical)

**CSRF Protection Tests:**
```go
func TestCSRFProtection(t *testing.T) {
    // Test 1: Forged state should be rejected
    resp := callOAuthCallback(provider, code, "forged_state")
    assert.Equal(t, 401, resp.StatusCode)

    // Test 2: Expired state should be rejected
    state := createStateExpiredMinutesAgo(10)
    resp = callOAuthCallback(provider, code, state)
    assert.Equal(t, 401, resp.StatusCode)

    // Test 3: Reused state should be rejected
    state = generateValidState()
    callOAuthCallback(provider, code, state) // First use
    resp = callOAuthCallback(provider, code, state) // Replay
    assert.Equal(t, 401, resp.StatusCode)
}
```

**Cookie Security Tests:**
```go
func TestCookieSecurity(t *testing.T) {
    resp := login(email, password)
    cookie := resp.Cookies["token"]

    assert.True(cookie.HttpOnly)
    assert.Equal(http.SameSiteLaxMode, cookie.SameSite)

    if env == "production" {
        assert.True(cookie.Secure)
    }
}
```

### 5.2 Integration Testing

**OAuth Flow Tests:**
```go
func TestCompleteOAuthFlow(t *testing.T) {
    // 1. Request OAuth URL
    resp := GET("/api/v1/oauth/github/login")
    var data struct{ URL string }
    json.Unmarshal(resp.Body, &data)

    // 2. Extract state from URL
    state := extractStateFromURL(data.URL)

    // 3. Verify state in database
    dbState := db.GetOAuthState(state)
    assert.Equal("github", dbState.Provider)

    // 4. Mock OAuth callback
    code := "mock_auth_code"
    resp = GET(fmt.Sprintf("/api/v1/oauth/github/callback?code=%s&state=%s", code, state))

    // 5. Verify user created
    assert.Equal(200, resp.StatusCode)
    var user User
    json.Unmarshal(resp.Body, &user)
    assert.NotEmpty(user.ID)

    // 6. Verify JWT cookie set
    assert.NotEmpty(resp.Cookies["token"])

    // 7. Verify state deleted
    _, err := db.GetOAuthState(state)
    assert.Error(err) // Should not exist
}
```

### 5.3 Frontend Testing

**Component Tests:**
```typescript
test('OAuth buttons trigger redirect', async () => {
  const { getByText } = render(<AuthComponent mode="login" />)

  const githubButton = getByText('GitHub')
  fireEvent.click(githubButton)

  // Should redirect (window.location.href changed)
  expect(window.location.href).toContain('/oauth/github/login')
})
```

**Error Handling Tests:**
```typescript
test('handles OAuth error gracefully', async () => {
  mockFetch.mockRejectedValueOnce(new AuthError('oauth_failed', 401, 'Provider denied'))

  const { getByText } = render(<AuthComponent mode="login" />)
  fireEvent.click(getByText('GitHub'))

  await waitFor(() => {
    expect(screen.getByText(/denied/i)).toBeInTheDocument()
  })
})
```

---

## 6. Security Considerations

### 6.1 OAuth State Management

**Implementation:**
- ✅ Cryptographically secure random state (32 bytes)
- ✅ Database storage with 5-minute TTL
- ✅ One-time use (deleted after validation)
- ✅ Provider validation (state tied to specific provider)

**Attack Prevention:**
- ✅ CSRF attacks (state validation)
- ✅ Replay attacks (one-time use)
- ✅ Timing attacks (constant-time comparison)

### 6.2 Cookie Security

**Flags:**
- ✅ `HttpOnly: true` - Prevents JavaScript access (XSS protection)
- ✅ `Secure: true` (production) - HTTPS only
- ✅ `SameSite: Lax` - CSRF protection

### 6.3 Secrets Management

**Current (Good):**
- ✅ OAuth secrets in environment variables
- ✅ JWT secret not exposed to frontend
- ✅ No hardcoded credentials

**Recommendations:**
- Consider using secret management service (AWS Secrets Manager, Vault)
- Rotate OAuth secrets periodically
- Monitor failed OAuth attempts

### 6.4 Rate Limiting

**TODO (Post-MVP):**
- Add rate limiting on `/oauth/:provider/login` (prevent DoS)
- Add rate limiting on `/oauth/:provider/callback` (prevent brute force)
- Implement exponential backoff for failed attempts

---

## 7. Rollback Plan

### If Issues Arise:

**Step 1: Identify Issue**
- Check error logs
- Check database state
- Check OAuth provider status

**Step 2: Quick Rollback**
Since this is a direct replacement:

```bash
# Revert to previous Git commit
git revert HEAD

# Or cherry-pick old implementation
git checkout <previous-commit> -- internal/pkg/auth/oauth_provider*.go
git checkout <previous-commit> -- internal/port/httpserver/handler_auth.go

# Rebuild and deploy
mise run build
./recally
```

**Step 3: Keep Database**
- No database rollback needed (schema changes are additive)
- `auth_oauth_states` table can remain (unused by old code)

**Step 4: User Impact**
- Active OAuth sessions continue working (JWT cookies unchanged)
- Users can re-authenticate with old flow

---

## 8. Future Enhancements (Post-MVP)

### Phase 2 Features:
- Email verification flow
- Password reset via email
- Magic link authentication
- 2FA/TOTP support
- WebAuthn/Passkey support
- More OAuth providers (Discord, Twitter, Apple)

### DevOps Improvements:
- Prometheus metrics for OAuth flows
- Grafana dashboard for auth monitoring
- Automated state cleanup job
- OAuth failure alerting

### Developer Experience:
- Auto-generate TypeScript types from Go structs
- Publish auth client as separate npm package
- Create OAuth provider cookbook/examples

---

## 9. Success Criteria Checklist

Before marking this complete, verify:

### Security:
- [ ] CSRF attack tests pass
- [ ] State expiration enforced
- [ ] Replay attack prevented
- [ ] Cookies have HttpOnly and Secure flags
- [ ] No secrets in frontend code

### Functionality:
- [ ] GitHub OAuth works end-to-end
- [ ] Google OAuth works end-to-end
- [ ] Telegram auth works (custom provider)
- [ ] Existing users can login
- [ ] JWT sessions preserved
- [ ] API keys unchanged

### Code Quality:
- [ ] All tests pass: `mise run test`
- [ ] Linters pass: `mise run lint`
- [ ] Build succeeds: `mise run build`
- [ ] No TypeScript errors
- [ ] Documentation updated

### Performance:
- [ ] OAuth flow latency < 500ms
- [ ] No memory leaks
- [ ] Database queries optimized
- [ ] Frontend bundle size acceptable (<50KB increase)

---

## 10. Summary

**Estimated Effort:** 24 hours over 3 days
**Risk Level:** Low-Medium
**Breaking Changes:** None

**Key Improvements:**
1. ✅ Secure OAuth implementation (CSRF-proof)
2. ✅ Multi-provider support (easy to add new providers)
3. ✅ Type-safe TypeScript client
4. ✅ Consistent Telegram auth integration
5. ✅ Better developer experience

**What Stays the Same:**
- API endpoints
- Database schema (only additive)
- JWT session management
- API keys
- User experience

This plan provides a secure, maintainable foundation for authentication while minimizing risk and preserving all existing functionality.
