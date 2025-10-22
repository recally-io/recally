# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Recally is an AI-powered memory assistant for digital content. It's a full-stack web application with:
- Go backend using Echo framework
- React frontend with TypeScript
- PostgreSQL database with ParadeDB extensions for full-text search
- Telegram bot integration
- Browser extensions support

## Essential Commands

### Development
```bash
# Full application with hot reload
mise run dev:backend  # Backend (go run, hot reload)
mise run run:ui       # Frontend dev server

# Production-like run (with build)
mise run run          # Build and run
mise run run:go       # Build and run with DEBUG_UI=true

# Database only
mise run db:up
```

### Code Quality
```bash
# Lint everything (Go + UI)
mise run lint

# Run tests
mise run test

# Generate code (SQL, Swagger)
mise run generate
```

### Database Management
```bash
# Create new migration
mise run migrate:new name=migration_name

# Apply migrations
mise run migrate:up

# Check migration status
mise run migrate:status

# Validate migrations
mise run migrate:validate

# Rollback last migration
mise run migrate:down

# Access PostgreSQL console
mise run psql
```

### Building
```bash
# Build everything
mise run build

# Build Docker image
mise run docker:build

# Run with Docker Compose
mise run docker:up
```

### Tool Management
```bash
# Install all tools defined in mise.toml
mise install

# List installed tools
mise list

# See all available tasks
mise tasks

# Check environment health
mise run doctor
```

## Architecture Overview

### Backend Structure (`/internal/`)
- **core/**: Business logic
  - `assistants/`: AI assistant functionality with conversation management
  - `bookmarks/`: Content processing, embedding, and search
  - `files/`: File storage and retrieval
  - `queue/`: Background jobs (crawling, embedding, summarization)

- **pkg/**: Shared packages
  - `auth/`: JWT, OAuth, API key authentication
  - `cache/`: Two-tier caching (DB + memory)
  - `db/`: Database layer using SQLC for type-safe queries
  - `llms/`: LLM integrations (OpenAI, Ollama)
  - `rag/`: RAG implementation for semantic search
  - `webreader/`: Web scraping and content extraction

- **port/**: External interfaces
  - `httpserver/`: REST API with Echo framework
  - `bots/`: Telegram bot implementation

### Frontend Structure (`/web/`)
- Uses React 18 with TypeScript
- TanStack Router for routing
- SWR for data fetching
- Tailwind CSS for styling
- PWA support with service workers

### Key Technologies
- **Database**: PostgreSQL with ParadeDB extensions (pg_search, pgvector)
- **Background Jobs**: River queue system
- **Web Scraping**: go-rod for browser automation
- **API Documentation**: Auto-generated Swagger/OpenAPI
- **Code Generation**: SQLC for type-safe SQL

## Development Workflow

1. **Environment Setup**:
   - Copy `env.example` to `.env`
   - Set required variables (especially `JWT_SECRET`, `OPENAI_API_KEY`)
   - Database runs on port 15432

2. **Code Generation**:
   - Run `mise run generate` after modifying SQL queries
   - SQL queries in `/database/queries/` generate Go code via SQLC
   - API spec auto-generated from Echo routes

3. **Database Changes**:
   - Create migrations with `mise run migrate:new name=feature_name`
   - Migrations stored in `/database/migrations/`
   - Always test migrations with `mise run migrate:up` and `mise run migrate:down`

4. **Testing**:
   - Backend tests: `mise run test`
   - Integration tests use real PostgreSQL (via `mise run db:up`)

## Important Patterns

### API Structure
- RESTful endpoints under `/api/v1/`
- Authentication via JWT tokens or API keys
- Request/response models in `/internal/port/httpserver/handlers/`

### Background Jobs
- Queue system using River
- Jobs defined in `/internal/core/queue/`
- Handles: web crawling, content embedding, summarization

### Content Processing Flow
1. User saves URL/content → Creates bookmark
2. Queue job crawls and extracts content
3. Content gets embedded for semantic search
4. Optional: AI generates summary
5. Content searchable via full-text and vector search

### Authentication
- JWT for web sessions
- API keys for programmatic access
- OAuth support (GitHub, Google)
- Telegram bot authentication via chat ID

## Common Tasks

### Adding New API Endpoint
1. Define handler in `/internal/port/httpserver/handlers/`
2. Add route in `/internal/port/httpserver/routes.go`
3. Run `mise run generate:spec` to update Swagger docs

### Adding Database Query
1. Write SQL in `/database/queries/`
2. Run `mise run generate:sql`
3. Use generated code in your Go files

### Modifying Frontend
1. Components in `/web/src/components/`
2. Routes in `/web/src/routes/`
3. API client in `/web/src/lib/api/`
4. Run `mise run run:ui` for hot reload

## Configuration Notes
- Service FQDN required for OAuth callbacks and webhooks
- ParadeDB provides PostgreSQL with built-in full-text search
- Browser service (go-rod) required for web scraping
- S3-compatible storage optional for file uploads
