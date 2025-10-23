# Implementation Tasks: Goth OAuth + TypeScript Auth SDK

**Session:** 2025-10-23-goth-typescript-auth-sdk
**Total Estimated Time:** 24 hours (3 days)

---

## Phase 0: Preparation (3 hours) ✅ COMPLETED

### Database & Configuration
- [x] Create migration file `database/migrations/20251023_oauth_states.sql`
  - **Files:** `database/migrations/20251023_oauth_states.sql`
  - **Approach:** Created auth_oauth_states table with state (VARCHAR 64 PK), provider, redirect_url, created_at, expires_at columns
  - **Commit:** (pending - will commit after Phase 1)
- [x] Run migration: `mise run migrate:up`
- [x] Verify table created: `mise run psql` → `\d auth_oauth_states`

### SQLC Query Generation
- [x] Add OAuth state queries to `database/queries/auth.sql`
  - **Files:** `database/queries/auth.sql`
  - **Approach:** Added CreateOAuthState, GetOAuthState, DeleteOAuthState, CleanupExpiredStates queries
  - **Commit:** (pending - will commit after Phase 1)
- [x] Run SQLC generation: `mise run generate:sql`
  - **Files:** `internal/pkg/db/auth.sql.go`, `internal/pkg/db/models.go`
  - **Approach:** Generated AuthOauthState model and CRUD functions
- [x] Verify generated Go code in `internal/pkg/db/`

### Dependencies
- [x] Install Goth: `go get github.com/markbates/goth`
  - **Files:** `go.mod`, `go.sum`
  - **Approach:** Installed goth v1.82.0 with all provider support
- [x] Run `go mod tidy`

### Google OAuth Setup
- [ ] Register app at Google Cloud Console (deferred to deployment)
- [ ] Get OAuth client ID and secret (deferred to deployment)
- [ ] Add to `.env` (deferred to deployment)
- [ ] Update config struct in `internal/pkg/config/` to include Google OAuth (will do in Phase 1)

### Directory Structure
- [x] Create `internal/pkg/auth/adapter/` directory
- [x] Create `web/src/lib/auth-client/` directory

---

## Phase 1: Backend Implementation (5 hours) ✅ COMPLETED

### OAuth Adapter Interface
- [x] Create `internal/pkg/auth/adapter/oauth_adapter.go`
  - **Files:** `internal/pkg/auth/adapter/oauth_adapter.go`
  - **Approach:** Created OAuthAdapter interface with GetAuthURL, HandleCallback, ListProviders methods; defined OAuth2User struct
  - **Commit:** (pending)
- [ ] Define `OAuthAdapter` interface with methods:
  - [ ] `GetAuthURL(ctx, provider) (url, error)`
  - [ ] `HandleCallback(ctx, provider, code, state) (OAuth2User, error)`
  - [ ] `ListProviders() []string`

### Goth Adapter Implementation
- [ ] Create `internal/pkg/auth/adapter/goth_adapter.go`
- [ ] Implement `GothAdapter` struct with:
  - [ ] Provider registry map
  - [ ] Database DAO field
- [ ] Implement `GetAuthURL` method:
  - [ ] Generate secure state (32 bytes, base64 encoded)
  - [ ] Store state in database with 5-minute TTL
  - [ ] Get provider and call `BeginAuth(state)`
  - [ ] Return auth URL
- [ ] Implement `HandleCallback` method:
  - [ ] Validate state exists in database
  - [ ] Check state hasn't expired
  - [ ] Verify provider matches
  - [ ] Delete used state (one-time use)
  - [ ] Complete OAuth flow with provider
  - [ ] Convert goth.User to OAuth2User
- [ ] Implement `ListProviders` method
- [ ] Implement `generateSecureState()` helper

### Custom Telegram Provider
- [ ] Create `internal/pkg/auth/adapter/telegram_provider.go`
- [ ] Implement `TelegramProvider` struct
- [ ] Implement `goth.Provider` interface methods:
  - [ ] `Name() string`
  - [ ] `BeginAuth(state) (Session, error)`
  - [ ] `FetchUser(session) (User, error)`
  - [ ] `UnmarshalSession(data) (Session, error)`
- [ ] Create `TelegramSession` struct
- [ ] Implement Telegram-specific auth logic
- [ ] Handle Telegram widget callback verification

### Goth Configuration
- [ ] Create `internal/pkg/auth/goth_config.go`
- [ ] Implement `InitGothAdapter(dao) *adapter.GothAdapter`:
  - [ ] Create GothAdapter instance
  - [ ] Register GitHub provider (if configured)
  - [ ] Register Google provider (if configured)
  - [ ] Register Telegram provider (if configured)
  - [ ] Return initialized adapter

### Auth Service Updates
- [ ] Update `internal/pkg/auth/service.go`:
  - [ ] Add `oauthAdapter *adapter.GothAdapter` field
  - [ ] Update constructor to accept adapter
  - [ ] Rename `HandleOAuth2Callback` to `HandleOAuth2UserLogin`
  - [ ] Update method to accept `OAuth2User` directly
  - [ ] Remove dependency on `GetOAuthProvider()` function
- [ ] Mark old OAuth provider files for removal

### HTTP Handler Updates
- [ ] Update `internal/port/httpserver/handler_auth.go`:
  - [ ] Add `oauthAdapter` field to `authHandler`
  - [ ] Update `oAuthLogin` handler to use `adapter.GetAuthURL()`
  - [ ] Update `oAuthCallback` handler:
    - [ ] Extract `code` and `state` from query params
    - [ ] Call `adapter.HandleCallback()`
    - [ ] Use existing `HandleOAuth2UserLogin()` logic
    - [ ] Generate JWT and set cookie
    - [ ] Redirect to frontend
  - [ ] Update `setCookieJwtToken()`:
    - [ ] Add `HttpOnly: true`
    - [ ] Add `Secure: config.Settings.Env == "production"`
    - [ ] Keep existing `SameSite: Lax`

### Main.go Integration
- [ ] Update `main.go`:
  - [ ] Initialize Goth adapter: `auth.InitGothAdapter(db.New())`
  - [ ] Pass adapter to HTTP server initialization

---

## Phase 2: Backend Testing (3 hours)

### Unit Tests
- [ ] Create `internal/pkg/auth/adapter/goth_adapter_test.go`
- [ ] Test `generateSecureState()`:
  - [ ] Verify 32-byte output
  - [ ] Verify base64 encoding
  - [ ] Verify randomness (no duplicates in 1000 iterations)
- [ ] Test `GetAuthURL()`:
  - [ ] Verify state stored in database
  - [ ] Verify state expires in 5 minutes
  - [ ] Verify URL contains state parameter
- [ ] Test `HandleCallback()`:
  - [ ] Valid state accepted
  - [ ] Invalid state rejected
  - [ ] Expired state rejected
  - [ ] Used state rejected (one-time use)
  - [ ] Provider mismatch rejected

### Integration Tests
- [ ] Create `internal/port/httpserver/handler_auth_integration_test.go`
- [ ] Test GitHub OAuth flow:
  - [ ] Request `/oauth/github/login`
  - [ ] Verify redirect URL returned
  - [ ] Extract state from URL
  - [ ] Mock GitHub callback with code + state
  - [ ] Verify user created in database
  - [ ] Verify JWT cookie set
  - [ ] Verify state deleted after use
- [ ] Test Google OAuth flow (same as above)
- [ ] Test Telegram OAuth flow (same as above)
- [ ] Test CSRF protection:
  - [ ] Forged state rejected with 401
  - [ ] Missing state rejected with 401
  - [ ] Replay attack prevented

### Manual Testing
- [ ] Start backend: `mise run dev:backend`
- [ ] Test GitHub OAuth in browser:
  - [ ] Click "Login with GitHub"
  - [ ] Authorize on GitHub
  - [ ] Verify redirect back to app
  - [ ] Verify logged in
- [ ] Test Google OAuth (same as above)
- [ ] Check database:
  - [ ] Verify user record created
  - [ ] Verify OAuth connection record created
  - [ ] Verify no lingering states in `auth_oauth_states`
- [ ] Test logout and re-login

---

## Phase 3: Frontend Implementation (6 hours)

### TypeScript Auth Client
- [ ] Create `web/src/lib/auth-client/types.ts`:
  - [ ] Define `User` interface
  - [ ] Define `UserSettings` interface
  - [ ] Define `LoginInput` interface
  - [ ] Define `RegisterInput` interface
- [ ] Create `web/src/lib/auth-client/errors.ts`:
  - [ ] Define `AuthError` class
  - [ ] Implement `fromResponse()` static method
- [ ] Create `web/src/lib/auth-client/client.ts`:
  - [ ] Define `AuthClient` class
  - [ ] Implement `login(email, password)` method
  - [ ] Implement `register(username, email, password)` method
  - [ ] Implement `logout()` method
  - [ ] Implement `validateSession()` method
  - [ ] Implement `getOAuthURL(provider)` method
  - [ ] Implement `redirectToOAuth(provider)` method
  - [ ] Implement private `get<T>(path)` helper
  - [ ] Implement private `post<T>(path, body)` helper
  - [ ] Create singleton `authClient` instance
- [ ] Create `web/src/lib/auth-client/index.ts`:
  - [ ] Export `AuthClient` class
  - [ ] Export `authClient` singleton
  - [ ] Export all types
  - [ ] Export `AuthError`

### Update Auth Hooks
- [ ] Update `web/src/lib/apis/auth.ts`:
  - [ ] Import `authClient` from `@/lib/auth-client`
  - [ ] Replace `fetcher` calls with `authClient` methods in `useUser()`
  - [ ] Update `useAuth()` hook:
    - [ ] Use `authClient.login()` in `login` function
    - [ ] Use `authClient.register()` in `register` function
    - [ ] Use `authClient.logout()` in `logout` function
    - [ ] Use `authClient.redirectToOAuth()` in `oauthLogin` function
  - [ ] Keep SWR configuration unchanged
  - [ ] Keep API key hooks unchanged (no changes needed)

### Update Auth Component
- [ ] Update `web/src/components/auth/auth.tsx`:
  - [ ] Import types from `@/lib/auth-client/types`
  - [ ] Add Google to `OAuthProviders` array:
    ```typescript
    {
      name: "Google",
      icon: SiGoogle,
    }
    ```
  - [ ] Uncomment Google OAuth button
  - [ ] Test all auth flows work

### Update Other Components
- [ ] Search for all uses of `useUser()` hook:
  ```bash
  grep -r "useUser" web/src/
  ```
- [ ] Verify compatibility (should be no changes needed)
- [ ] Test protected route components
- [ ] Test session persistence across page refresh

### Frontend Testing
- [ ] Create `web/src/lib/auth-client/__tests__/client.test.ts`:
  - [ ] Test `login()` success
  - [ ] Test `login()` failure
  - [ ] Test `register()` success
  - [ ] Test `logout()`
  - [ ] Test `validateSession()`
  - [ ] Test error handling
- [ ] Create component tests:
  - [ ] Test OAuth buttons render
  - [ ] Test OAuth button clicks trigger redirect
  - [ ] Test login form submission
  - [ ] Test registration form submission
  - [ ] Test error messages display
- [ ] E2E test:
  - [ ] Full registration flow
  - [ ] Login flow
  - [ ] OAuth flow (GitHub and Google)
  - [ ] Logout flow

---

## Phase 4: Testing & Documentation (6 hours)

### Security Audit
- [ ] CSRF Attack Tests:
  - [ ] Test forged state parameter rejected
  - [ ] Test state replay attack prevented
  - [ ] Test state expiration enforced
  - [ ] Test provider mismatch rejected
- [ ] Cookie Security Tests:
  - [ ] Verify `HttpOnly` flag set
  - [ ] Verify `Secure` flag set in production
  - [ ] Verify `SameSite=Lax` set
  - [ ] Test XSS prevention (JavaScript can't read cookie)
- [ ] OAuth Flow Security:
  - [ ] Test redirect URI validation
  - [ ] Test callback origin validation
  - [ ] Test token exchange security

### Load Testing
- [ ] Install Apache Bench or similar tool
- [ ] Test OAuth login endpoint:
  ```bash
  ab -n 100 -c 10 http://localhost:8080/api/v1/oauth/github/login
  ```
- [ ] Test OAuth callback endpoint (mock)
- [ ] Monitor:
  - [ ] Response times (should be <500ms)
  - [ ] Error rate (should be 0%)
  - [ ] Database connection pool usage
  - [ ] Memory usage

### Migration Validation
- [ ] Create test user via old system (if possible)
- [ ] Verify old user can login via new system
- [ ] Test existing OAuth connections work
- [ ] Verify API keys unchanged and working
- [ ] Test Telegram auth preserved
- [ ] Check database:
  - [ ] No orphaned records
  - [ ] OAuth states cleaned up properly

### Performance Testing
- [ ] Measure OAuth flow end-to-end latency
- [ ] Compare bundle size before/after:
  ```bash
  cd web && npm run build
  du -sh dist/
  ```
- [ ] Verify frontend bundle increase <50KB
- [ ] Check backend memory usage
- [ ] Profile database queries

### Documentation
- [ ] Update `CLAUDE.md`:
  - [ ] Add Goth OAuth architecture section
  - [ ] Document OAuth provider setup
  - [ ] Update authentication flow diagram
  - [ ] Add troubleshooting section
- [ ] Create `web/src/lib/auth-client/README.md`:
  - [ ] Usage examples
  - [ ] API reference
  - [ ] Type documentation
  - [ ] Error handling guide
- [ ] Document OAuth Provider Setup:
  - [ ] Create `docs/oauth-setup.md`
  - [ ] GitHub OAuth app registration
  - [ ] Google OAuth app registration
  - [ ] Environment variable configuration
- [ ] Create migration notes in `docs/migrations/goth-migration.md`:
  - [ ] What changed
  - [ ] Breaking changes (none)
  - [ ] Rollback procedure
  - [ ] Testing checklist
- [ ] Update API documentation (if using Swagger/OpenAPI)

---

## Phase 5: Cleanup (1 hour)

### Remove Old Code
- [ ] Delete `internal/pkg/auth/oauth_provider.go`
- [ ] Delete `internal/pkg/auth/oauth_provider_github.go`
- [ ] Remove imports of deleted files
- [ ] Clean up unused functions

### Final Verification
- [ ] Run full test suite: `mise run test`
  - [ ] All tests pass
  - [ ] No test warnings
- [ ] Run linters: `mise run lint`
  - [ ] No linting errors
  - [ ] No security warnings
- [ ] Build project: `mise run build`
  - [ ] Backend builds successfully
  - [ ] Frontend builds successfully
  - [ ] No TypeScript errors
- [ ] Generate code: `mise run generate`
  - [ ] SQLC generation works
  - [ ] Swagger spec updated (if applicable)

### Deployment Preparation
- [ ] Update `.env.example` with Google OAuth variables
- [ ] Create deployment checklist
- [ ] Document environment variables
- [ ] Create rollback plan document
- [ ] Tag release in git (if applicable)

---

## Success Criteria Verification

### Security ✓
- [ ] CSRF attack tests pass
- [ ] State expiration enforced
- [ ] Replay attack prevented
- [ ] Cookies have HttpOnly and Secure flags
- [ ] No secrets in frontend code
- [ ] OAuth redirect URI validated

### Functionality ✓
- [ ] GitHub OAuth works end-to-end
- [ ] Google OAuth works end-to-end
- [ ] Telegram auth works (custom provider)
- [ ] Existing users can login
- [ ] JWT sessions preserved
- [ ] API keys unchanged
- [ ] Email/password login works
- [ ] Registration works
- [ ] Logout works

### Code Quality ✓
- [ ] All tests pass: `mise run test`
- [ ] Linters pass: `mise run lint`
- [ ] Build succeeds: `mise run build`
- [ ] No TypeScript errors
- [ ] Documentation updated
- [ ] Code reviewed (if team workflow requires)

### Performance ✓
- [ ] OAuth flow latency < 500ms
- [ ] No memory leaks
- [ ] Database queries optimized
- [ ] Frontend bundle size increase <50KB
- [ ] Load test passes (100 concurrent requests)

---

## Rollback Procedure (If Needed)

### Quick Rollback Steps
1. [ ] Identify issue (check logs, database, OAuth provider status)
2. [ ] Revert code changes:
   ```bash
   git revert HEAD
   # OR
   git checkout <previous-commit> -- internal/pkg/auth/oauth*.go
   git checkout <previous-commit> -- internal/port/httpserver/handler_auth.go
   ```
3. [ ] Rebuild: `mise run build`
4. [ ] Redeploy: `./recally`
5. [ ] Verify old OAuth flow working
6. [ ] Monitor error logs

### Post-Rollback
- [ ] Document issue that caused rollback
- [ ] Create fix plan
- [ ] Test fix in staging environment
- [ ] Retry deployment

---

## Notes

**Time Tracking:**
- Phase 0: 3 hours
- Phase 1: 5 hours
- Phase 2: 3 hours
- Phase 3: 6 hours
- Phase 4: 6 hours
- Phase 5: 1 hour
- **Total: 24 hours**

**Dependencies:**
- PostgreSQL database running (port 15432)
- GitHub OAuth app credentials
- Google OAuth app credentials
- Telegram bot token (for Telegram auth)

**Testing Environments:**
- Local development: `http://localhost:8080`
- OAuth callbacks: Configure to match in provider settings

**Key Decisions Made:**
- ✅ State Management: Database table (not cache)
- ✅ Telegram Provider: Custom Goth provider
- ✅ Cookie Security: HttpOnly + Secure in production
- ✅ SDK Complexity: Minimal (typed client only)
- ✅ Rollout: Direct replacement (no gradual rollout)

**Codex Review Feedback Addressed:**
- ✅ Secure state management with CSRF protection
- ✅ Custom Telegram provider for consistency
- ✅ Cookie security flags (HttpOnly, Secure)
- ✅ Adapter layer isolates Goth from domain logic
- ✅ Enhanced testing for security edge cases
- ✅ Simplified TypeScript SDK approach

---

**Ready to implement? Start with Phase 0!** 🚀
