# Goth Authentication Migration - Technical Design

## Technical Overview

### Architecture Approach
We will implement a clean Goth-based authentication system alongside the existing OAuth implementation, creating new database tables and handlers while maintaining backward compatibility. This approach allows for a gradual migration with fallback capabilities and zero-downtime deployment.

### Technology Stack
- **Backend**: Go with Echo framework (existing)
- **OAuth Library**: Goth (github.com/markbates/goth) v1.78+
- **Database**: PostgreSQL with new dedicated tables
- **Session Management**: Goth's built-in session handling with database persistence
- **Security**: CSRF protection, secure cookies, JWT integration

### Key Design Decisions
1. **Clean Schema**: New tables specifically designed for Goth's data structures
2. **Dual Implementation**: Run Goth alongside existing OAuth during migration
3. **Provider Registry**: Dynamic provider configuration through environment variables
4. **Session Persistence**: Database-backed sessions for scalability
5. **JWT Integration**: Maintain existing JWT token system for API compatibility

## System Architecture

### High-Level Components
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Frontend      │────│   Echo Router    │────│   Goth Handler  │
│   (React/Web)   │    │   Middleware     │    │   Registry      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │                        │
                       ┌────────▼────────┐              │
                       │  Auth Middleware │              │
                       │  (JWT/Session)   │              │
                       └────────┬────────┘              │
                                │                        │
                    ┌───────────▼───────────┐           │
                    │    Session Store      │           │
                    │   (Database-backed)   │           │
                    └───────────┬───────────┘           │
                                │                        │
                    ┌───────────▼───────────┐    ┌──────▼──────┐
                    │    PostgreSQL DB      │    │   OAuth     │
                    │   - goth_sessions     │    │  Providers  │
                    │   - goth_providers    │    │ (External)  │
                    │   - goth_users        │    └─────────────┘
                    └───────────────────────┘
```

### Data Flow
1. **Login Initiation**: User clicks OAuth provider → Goth redirects to provider
2. **OAuth Callback**: Provider redirects back → Goth handles callback
3. **User Resolution**: Goth extracts user info → Create/update user record
4. **Session Creation**: Generate session → Store in database + set cookies
5. **JWT Generation**: Create JWT token → Return to client for API access

### Component Interactions
- **Echo Router**: Routes OAuth requests to Goth handlers
- **Goth Registry**: Manages configured OAuth providers
- **Session Store**: Persists sessions in PostgreSQL for scalability
- **Auth Middleware**: Validates sessions and JWT tokens
- **User Service**: Manages user accounts and OAuth connections

## Data Design

### New Database Schema

#### goth_providers Table
```sql
CREATE TABLE goth_providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    client_id VARCHAR(255) NOT NULL,
    client_secret_encrypted TEXT NOT NULL,
    scopes TEXT[] DEFAULT '{}',
    auth_url VARCHAR(500),
    token_url VARCHAR(500),
    profile_url VARCHAR(500),
    callback_url VARCHAR(500) NOT NULL,
    enabled BOOLEAN DEFAULT true,
    config JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_goth_providers_name ON goth_providers(name);
CREATE INDEX idx_goth_providers_enabled ON goth_providers(enabled);
```

#### goth_users Table
```sql
CREATE TABLE goth_users (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider_name VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    name VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    nickname VARCHAR(255),
    avatar_url TEXT,
    access_token_encrypted TEXT,
    refresh_token_encrypted TEXT,
    expires_at TIMESTAMP WITH TIME ZONE,
    raw_data JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(provider_name, provider_user_id)
);

CREATE INDEX idx_goth_users_user_id ON goth_users(user_id);
CREATE INDEX idx_goth_users_provider ON goth_users(provider_name, provider_user_id);
CREATE INDEX idx_goth_users_email ON goth_users(email);
```

#### goth_sessions Table
```sql
CREATE TABLE goth_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    session_data JSONB NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_accessed TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_goth_sessions_user_id ON goth_sessions(user_id);
CREATE INDEX idx_goth_sessions_expires_at ON goth_sessions(expires_at);
```

### Data Relationships
- `goth_users.user_id` → `users.id` (existing users table)
- `goth_sessions.user_id` → `users.id` (session ownership)
- `goth_users.provider_name` → `goth_providers.name` (provider reference)

### Data Validation
- OAuth tokens encrypted at rest using AES-256-GCM
- Provider credentials encrypted in database
- Session data sanitized and validated
- Foreign key constraints ensure data integrity

### Migration Strategy
1. Create new tables alongside existing schema
2. Implement data sync between old and new systems during transition
3. Gradual migration of user sessions to new format
4. Cleanup of old tables after successful migration

## API Design

### OAuth Endpoints

#### Provider Discovery
```http
GET /api/v1/auth/providers
Response: {
  "providers": [
    {
      "name": "github",
      "display_name": "GitHub", 
      "auth_url": "/api/v1/auth/github",
      "enabled": true
    }
  ]
}
```

#### OAuth Initiation
```http
GET /api/v1/auth/{provider}
Query Parameters:
  - redirect_uri: Optional redirect after success
  - state: CSRF protection token
Response: HTTP 302 redirect to provider
```

#### OAuth Callback
```http
GET /api/v1/auth/{provider}/callback
Query Parameters:
  - code: Authorization code from provider
  - state: CSRF protection token
Response: {
  "success": true,
  "user": {
    "id": 123,
    "email": "user@example.com",
    "name": "John Doe"
  },
  "token": "jwt_token_here",
  "redirect_url": "/dashboard"
}
```

#### Session Management
```http
POST /api/v1/auth/logout
DELETE /api/v1/auth/sessions/{session_id}
GET /api/v1/auth/sessions
```

#### User Connection Management
```http
GET /api/v1/user/connections
Response: {
  "connections": [
    {
      "provider": "github",
      "email": "user@example.com",
      "name": "John Doe",
      "avatar_url": "https://.../avatar.png",
      "connected_at": "2023-10-27T10:00:00Z"
    }
  ]
}

DELETE /api/v1/user/connections/{provider}
Response: HTTP 204 No Content
```

### Authentication Flow
1. **GET /auth/{provider}**: Initiate OAuth with CSRF state
2. **Provider Authorization**: User authorizes on provider site
3. **GET /auth/{provider}/callback**: Handle callback and create session
4. **JWT Token**: Return JWT for API access + set session cookies

### Error Handling
- **400 Bad Request**: Invalid provider or malformed request
- **401 Unauthorized**: OAuth authorization failed
- **403 Forbidden**: CSRF token mismatch
- **500 Internal Error**: Provider communication failure

## Security Considerations

### Authentication Mechanisms
- **OAuth 2.0**: Standard OAuth flows with PKCE where supported
- **CSRF Protection**: State parameter validation for all OAuth flows
- **Session Security**: HTTP-only, secure, SameSite cookies
- **JWT Integration**: Maintain existing JWT system for API access

### Data Protection
- **Token Encryption**: All OAuth tokens encrypted at rest (AES-256-GCM)
- **Credential Storage**: Provider secrets encrypted in database
- **Key Management**: Use existing key management from current auth system
- **Data Minimization**: Store only necessary user data from providers

### Input Validation
- **OAuth Parameters**: Validate all callback parameters
- **User Data**: Sanitize all user information from providers
- **Session Data**: Validate session structure and expiration
- **Provider Config**: Validate provider configuration at startup

### Access Control
- **Provider Isolation**: Each provider isolated in separate user records
- **Session Scoping**: Sessions tied to specific users and expire appropriately
- **Admin Controls**: Provider management restricted to admin users
- **Audit Logging**: All authentication events logged with user context

## Performance & Scalability

### Performance Targets
- **OAuth Flow**: Complete within 3 seconds (requirement)
- **Session Lookup**: < 50ms for session validation
- **Provider Config**: Cached in memory, refreshed every 5 minutes
- **Database Queries**: Optimized with proper indexing

### Caching Strategies
- **Provider Registry**: In-memory cache of provider configurations
- **Session Cache**: Optional Redis cache for high-traffic scenarios
- **User Data**: Cache user OAuth data for 15 minutes
- **CSRF Tokens**: In-memory store with automatic cleanup

### Database Optimization
- **Indexing**: Strategic indexes on lookup columns
- **Connection Pooling**: Use existing database pool configuration
- **Query Optimization**: Minimize N+1 queries in user resolution
- **Cleanup Jobs**: Automatic cleanup of expired sessions

### Scaling Considerations
- **Stateless Design**: All handlers stateless for horizontal scaling
- **Database Sessions**: Support multiple application instances
- **Provider Failover**: Graceful handling of provider downtime
- **Rate Limiting**: Implement rate limiting on OAuth endpoints

## Implementation Approach

### Development Phases

#### Phase 1: Foundation (Week 1)
- Install and configure Goth library
- Create new database tables and migrations
- Implement provider registry and configuration loader
- Set up encrypted credential storage

#### Phase 2: Core Integration (Week 2)
- Implement Goth handlers for OAuth flows
- Create database-backed session store
- Integrate with Echo routing and middleware
- Implement CSRF protection

#### Phase 3: User Management (Week 3)
- User creation and linking logic
- JWT token integration
- Session management endpoints
- Provider connection management

#### Phase 4: Testing & Security (Week 4)
- Comprehensive OAuth flow testing
- Security validation and penetration testing
- Performance benchmarking
- Error handling and edge cases

#### Phase 5: Migration & Deployment (Week 5)
- Data migration scripts from existing OAuth
- Feature flag implementation
- Gradual rollout procedures
- Monitoring and alerting setup

### Testing Strategy
- **Unit Tests**: All handler and service logic
- **Integration Tests**: Full OAuth flows with test providers
- **Security Tests**: CSRF, token validation, session security
- **Performance Tests**: Load testing of OAuth endpoints
- **Migration Tests**: Data migration validation

### Deployment Plan
1. **Staging Deployment**: Full Goth system in staging environment
2. **Feature Flag**: Deploy to production behind feature flag
3. **Gradual Rollout**: Enable for percentage of users
4. **Migration Window**: Migrate existing users during low-traffic period
5. **Full Activation**: Complete switch to Goth system
6. **Cleanup**: Remove old OAuth code after validation period

## Migration from Existing System

### Data Migration Strategy
1. **Parallel Operation**: Run both systems during transition
2. **User Mapping**: Map existing `auth_user_oauth_connections` to `goth_users`
3. **Session Migration**: Gradual migration of active sessions
4. **Token Migration**: Encrypted migration of OAuth tokens
5. **Validation**: Verify all data integrity during migration

### Compatibility Layer
- Maintain existing API endpoints during transition
- Implement response format adapters for backward compatibility
- Support both old and new session formats
- Gradual deprecation of legacy endpoints

### Rollback Procedures
- Feature flag for instant rollback
- Database snapshot before migration
- Automated rollback scripts
- Health check monitoring during migration

### Success Metrics
- Zero data loss during migration
- All existing users can authenticate
- API response times within targets
- No increase in authentication errors

This technical design addresses all EARS requirements while providing a robust, secure, and scalable OAuth implementation using Goth. The clean database schema approach ensures optimal performance and maintainability going forward.