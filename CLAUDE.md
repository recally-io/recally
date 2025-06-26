# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Important

1. When searching or analyzing code structure or syntax, always use `ast-grep --lang '<language>' -p '<pattern>'` by default. Only use plain-text search tools like `rg` or `grep` if explicitly instructed to do so. Clearly document your search patterns and findings in the session log.
2. Create a session log file in the format `{time}-{name}.log` (e.g., `20240612-implement-auth.log`) in the `.claude/sessions/` directory. The filename should include the current timestamp and a descriptive session title.
3. For any task involving multiple steps, always decompose it into clear, manageable sub-tasks. Develop a detailed, step-by-step plan and checklist, and document all progress in a dedicated Markdown session log. Update the log as you complete each sub-task, ensuring traceability and transparency. If a task can be parallelized, create and coordinate sub-agents as needed to efficiently complete sub-tasks. Clearly document the responsibilities and outputs of each sub-agent in the session log.
4. Frequently commit changes to git using clear, conventional commit messages as checkpoints. This ensures that progress is well-documented and allows for easy rollback if needed.
5. Prioritize clarity and explicitness in all reasoning and documentation. For every instruction, proceed step by step, explaining your thought process and actions. Never make assumptions—ask for clarification if any requirement is ambiguous or incomplete.

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
make run

# Backend only (with DEBUG_UI=true)
make run-go

# Frontend only
make run-ui

# Database only
make db-up
```

### Code Quality
```bash
# Lint everything (Go + UI)
make lint

# Run tests
make test

# Generate code (SQL, Swagger, Go bindata)
make generate
```

### Database Management
```bash
# Create new migration
make migrate-new name=migration_name

# Apply migrations
make migrate-up

# Revert migrations
make migrate-down

# Access PostgreSQL console
make psql
```

### Building
```bash
# Build everything
make build

# Build Docker image
make docker-build

# Run with Docker Compose
make docker-up
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
   - Run `make generate` after modifying SQL queries
   - SQL queries in `/database/queries/` generate Go code via SQLC
   - API spec auto-generated from Echo routes

3. **Database Changes**:
   - Create migrations with `make migrate-new name=feature_name`
   - Migrations stored in `/database/migrations/`
   - Always test migrations with `make migrate-up` and `make migrate-down`

4. **Testing**:
   - Backend tests: `make test`
   - Integration tests use real PostgreSQL (via `make db-up`)

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
3. Run `make generate-spec` to update Swagger docs

### Adding Database Query
1. Write SQL in `/database/queries/`
2. Run `make generate-sql`
3. Use generated code in your Go files

### Modifying Frontend
1. Components in `/web/src/components/`
2. Routes in `/web/src/routes/`
3. API client in `/web/src/lib/api/`
4. Run `make run-ui` for hot reload

## Configuration Notes
- Service FQDN required for OAuth callbacks and webhooks
- ParadeDB provides PostgreSQL with built-in full-text search
- Browser service (go-rod) required for web scraping
- S3-compatible storage optional for file uploads
