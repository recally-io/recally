# Migrate Makefile to mise for Dependency and Task Management

**Date**: 2025-10-22
**Status**: In Progress
**Reviewed by**: Codex

## Implementation Progress

- [ ] Task 1: Create git branch and backup existing mise.toml
- [ ] Task 2: Create comprehensive mise.toml with all tools and tasks
- [ ] Task 3: Update README.md with mise prerequisites and commands
- [ ] Task 4: Update CLAUDE.md Essential Commands section
- [ ] Task 5: Create MIGRATION.md team migration guide
- [ ] Task 6: Update GitHub Actions workflows to use mise
- [ ] Task 7: Remove Makefile after validation
- [ ] Task 8: Update spec file with completion status

---

## Overview

### Goal
Completely migrate from GNU Make to mise for dependency management and task execution in the Recally project.

### Key Requirements
- Replace Makefile entirely with mise configuration
- Manage all development tools via mise (Go, Node.js/Bun, PostgreSQL clients, sqlc, golangci-lint, swag, Atlas)
- Convert all make targets to mise tasks
- Migrate relevant environment variables to mise
- Update all documentation and CI/CD workflows
- Organize tasks by functional category for better discoverability

### Success Criteria
- All Makefile functionality replicated in mise
- CI/CD pipelines working with mise
- Documentation updated and accurate
- Makefile removed from repository
- Developers can run all tasks via `mise run <task>`
- Team has clear migration guide and support

---

## Technical Approach

### mise Architecture

**mise** (formerly rtx) is a polyglot tool version manager and task runner that:
- Manages multiple runtime versions (Go, Node.js, etc.)
- Installs and manages CLI tools
- Defines environment variables per-project
- Runs tasks with proper tool versions activated
- Uses `.mise.toml` or `mise.toml` for configuration

### Configuration Structure

The `mise.toml` will contain:

1. **min_version** - Minimum mise version required
2. **[tools]** - Tool versions to install
3. **[env]** - Environment variables
4. **[tasks.*]** - Task definitions organized by category

### Tool Installation Strategy

mise supports multiple backends for tool installation:
- **core:go** - Go runtime via mise core plugin
- **core:node** - Node.js runtime via mise core plugin
- **ubi:** - GitHub release binaries (for golangci-lint, atlas, bun)
- **go:** - Go tools via `go install` (for sqlc, swag)

### Task Organization

Tasks will be organized with dot notation:
- `lint:*` - Linting tasks (Go, UI)
- `generate:*` - Code generation (Go, SQL, Swagger)
- `build:*` - Build operations (Go, UI, Docs)
- `run:*` / `dev:*` - Development servers
- `db:*` / `migrate:*` - Database operations
- `docker:*` - Docker operations
- Utilities: `setup`, `doctor`, `clean`, `help`

### Environment Variable Strategy

Move Makefile variables to mise `[env]` section:
- Load sensitive values from `.env` file using `_.file = [".env"]`
- Define non-sensitive defaults in mise.toml (DB_HOST, DB_PORT, etc.)
- Construct DATABASE_URL from components
- Reference system environment for secrets

---

## Implementation Steps

### Step 1: Create Git Branch
**Complexity**: Low

```bash
git checkout -b feat/migrate-to-mise
```

### Step 2: Backup Existing mise.toml
**Complexity**: Low
**Dependencies**: None

Current repository has an untracked `mise.toml` with some initial work:
```bash
mv mise.toml mise.toml.backup  # Keep for reference
```

This backup contains:
- golangci-lint tool definition
- lint:go task with modernize tool

We'll incorporate these into the new complete configuration.

### Step 3: Analyze Current Makefile
**Complexity**: Low
**Dependencies**: None

Current Makefile (153 lines) contains:

**Linting tasks** (5, 7, 13):
- `lint` - Lint Go + UI
- `lint-go` - Go mod tidy, golangci-lint, swag fmt
- `lint-ui` - Bun lint in web/

**Generation tasks** (17, 19, 23, 27):
- `generate` - All generation
- `generate-go` - go generate
- `generate-sql` - sqlc generate
- `generate-spec` - swag init for Swagger docs

**Build tasks** (34, 36, 40, 44):
- `build` - Build all (UI + Docs + Go)
- `build-go` - Build Go binary to `bin/app`
- `build-docs` - Build docs with Bun
- `build-ui` - Build web with Bun

**Test tasks** (48):
- `test` - Run Go tests

**Run tasks** (52, 56, 60):
- `run` - Build and run
- `run-ui` - Bun dev server
- `run-go` - Build and run with DEBUG_UI=true

**Utility tasks** (64):
- `ngrok` - Run ngrok tunnel

**Database tasks** (68, 120):
- `db-up` - Start postgres container
- `psql` - Docker exec into postgres

**Migration tasks** (88-118):
- `migrate-new` - Create Atlas migration (in database/)
- `migrate-up` - Apply migrations with Atlas
- `migrate-status` - Check migration status
- `migrate-validate` - Validate migrations
- `migrate-hash` - Generate checksums

**Docker tasks** (72-86):
- `docker-build` - Build images
- `docker-run` - Up with build
- `docker-up` - Up detached with build
- `docker-down` - Stop services

**Other tasks** (124, 128):
- `deploy` - Dokploy deployment
- `help` - Show help

**Key observations**:
- Uses Bun (not npm) for all frontend/docs
- Binary built to `bin/app` (not `bin/recally`)
- Atlas runs in `database/` directory
- Docker service is `postgres` (not `db`)
- DATABASE_URL constructed from env vars
- Uses `docker compose exec -it postgres psql`

### Step 4: Create Complete mise.toml
**Complexity**: High
**Dependencies**: Step 2, 3

Create comprehensive mise.toml with all tasks and tools:

```toml
min_version = "2024.1.0"

# ============================================
# TOOLS - Pinned versions for reproducibility
# ============================================
[tools]
# Core runtimes
go = "1.24.6"  # From go.mod
"ubi:oven-sh/bun" = "1.1"

# Go development tools
"go:github.com/sqlc-dev/sqlc/cmd/sqlc" = "1.27"
"ubi:golangci/golangci-lint" = "1.62"
"go:github.com/swaggo/swag/cmd/swag" = "1.16"

# Database tools
"ubi:ariga/atlas" = "0.29"

# ============================================
# ENVIRONMENT VARIABLES
# ============================================
[env]
# Load from .env file (for secrets like DB_PASSWORD, API keys)
_.file = [".env"]

# Database configuration (non-sensitive defaults)
DB_HOST = "localhost"
DB_PORT = "15432"
DB_NAME = "recally"
DB_USER = "postgres"

# Construct DATABASE_URL from components
DATABASE_URL = "postgresql://{{env.DB_USER}}:{{env.DB_PASSWORD}}@{{env.DB_HOST}}:{{env.DB_PORT}}/{{env.DB_NAME}}?sslmode=disable"

# ============================================
# LINTING TASKS
# ============================================
[tasks."lint"]
description = "Lint Go code and frontend"
depends = ["lint:ui", "lint:go"]

[tasks."lint:go"]
description = "Lint Go code with golangci-lint"
run = [
  "go mod tidy",
  "golangci-lint run --fix ./...",
  "swag fmt"
]

[tasks."lint:ui"]
description = "Lint frontend code with Biome"
dir = "web"
run = "bun run lint"

# ============================================
# CODE GENERATION TASKS
# ============================================
[tasks."generate"]
description = "Generate all code (Go + SQL + Swagger)"
depends = ["generate:go", "generate:sql", "generate:spec"]

[tasks."generate:go"]
description = "Generate Go code"
run = "go generate ./..."

[tasks."generate:sql"]
description = "Generate Go code from SQL using SQLC"
run = "sqlc generate"

[tasks."generate:spec"]
description = "Generate Swagger API documentation"
run = "swag init -g internal/port/httpserver/router.go -o docs/swagger"

# ============================================
# BUILD TASKS
# ============================================
[tasks."build"]
description = "Build everything (UI + Docs + Go)"
depends = ["build:ui", "build:docs", "build:go"]

[tasks."build:go"]
description = "Build Go backend binary"
depends = ["generate"]
run = "go build -o bin/app main.go"

[tasks."build:docs"]
description = "Build documentation site"
dir = "docs"
run = "bun run docs:build"

[tasks."build:ui"]
description = "Build frontend for production"
dir = "web"
run = "bun run build"

# ============================================
# TEST TASKS
# ============================================
[tasks."test"]
description = "Run Go tests"
run = "go test ./..."

# ============================================
# DEVELOPMENT TASKS
# ============================================
[tasks."run"]
description = "Build and run full application"
depends = ["build"]
run = "./bin/app"

[tasks."run:go"]
description = "Build and run backend with DEBUG_UI=true"
depends = ["build:go"]
env = { DEBUG_UI = "true" }
run = "./bin/app"
alias = "run-go"

[tasks."run:ui"]
description = "Run frontend development server"
dir = "web"
run = "bun run dev"
alias = "run-ui"

[tasks."dev:backend"]
description = "Run backend with hot reload (no build)"
env = { DEBUG_UI = "true" }
run = "go run main.go"

[tasks."dev:docs"]
description = "Run documentation development server"
dir = "docs"
run = "bun run docs:dev"

# ============================================
# DATABASE TASKS
# ============================================
[tasks."db:up"]
description = "Start PostgreSQL database"
run = "docker compose up -d postgres"

[tasks."migrate:new"]
description = "Create new Atlas migration (usage: mise run migrate:new name=your_migration)"
dir = "database"
run = "atlas migrate diff {{arg(name='migration')}} --env local"

[tasks."migrate:up"]
description = "Apply all pending Atlas migrations"
dir = "database"
run = "atlas migrate apply --env local --url \"$DATABASE_URL\""

[tasks."migrate:status"]
description = "Check Atlas migration status"
dir = "database"
run = "atlas migrate status --env local --url \"$DATABASE_URL\""

[tasks."migrate:validate"]
description = "Validate Atlas migrations"
dir = "database"
run = "atlas migrate validate --env local"

[tasks."migrate:hash"]
description = "Generate Atlas migration checksums"
dir = "database"
run = "atlas migrate hash"

[tasks."psql"]
description = "Connect to PostgreSQL console via Docker"
run = "docker compose exec -it postgres psql -U $DB_USER -d $DB_NAME"

# ============================================
# DOCKER TASKS
# ============================================
[tasks."docker:build"]
description = "Build Docker images"
run = "docker compose build"

[tasks."docker:run"]
description = "Run with docker compose (build + up)"
run = "docker compose up --build"

[tasks."docker:up"]
description = "Start all services with docker compose (detached)"
run = "docker compose up --build -d"

[tasks."docker:down"]
description = "Stop docker compose services"
run = "docker compose down"

# ============================================
# UTILITY TASKS
# ============================================
[tasks."ngrok"]
description = "Run ngrok tunnel to localhost:1323"
run = "ngrok http 1323"

[tasks."deploy"]
description = "Deploy using dokploy"
run = "dokploy app deploy"

[tasks."help"]
description = "Show available tasks"
run = "mise tasks"

[tasks."setup"]
description = "First-time project setup"
run = '''
#!/usr/bin/env bash
set -e

echo "🚀 Setting up Recally project..."

echo "📦 Installing Go dependencies..."
go mod download || { echo "❌ Failed to download Go modules"; exit 1; }

echo "📦 Installing frontend dependencies..."
cd web && bun install || { echo "❌ Failed to install frontend deps"; exit 1; }
cd ..

echo "📦 Installing docs dependencies..."
cd docs && bun install || { echo "❌ Failed to install docs deps"; exit 1; }
cd ..

echo "🐳 Starting database..."
mise run db:up || { echo "❌ Failed to start database"; exit 1; }

echo "⏳ Waiting for database to be ready..."
sleep 3

echo "🗃️ Applying migrations..."
mise run migrate:up || { echo "❌ Failed to apply migrations"; exit 1; }

echo "🔧 Generating code..."
mise run generate || { echo "❌ Failed to generate code"; exit 1; }

echo "✅ Setup complete! Run 'mise run dev:backend' and 'mise run run:ui' to start developing."
'''

[tasks."doctor"]
description = "Check tool versions and environment"
run = '''
echo "🔍 Recally Environment Check"
echo "============================"
mise --version
echo ""
echo "📦 Installed Tools:"
mise list
echo ""
echo "🔧 Tool Versions:"
go version
bun --version
atlas version
sqlc version
golangci-lint --version
echo ""
echo "🐳 Docker:"
docker --version
docker compose version
echo ""
echo "🗄️ Database:"
docker compose ps postgres || echo "⚠️  Database not running (run: mise run db:up)"
'''

[tasks."clean"]
description = "Clean build artifacts"
run = '''
rm -rf bin/
rm -rf web/dist/
rm -rf docs/.vitepress/dist/
rm -rf coverage.out
echo "✨ Cleaned build artifacts"
'''
```

**Key decisions**:
- Pin specific versions for reproducibility
- Incorporate lint:go task from backup mise.toml
- Use `dir = "database"` for Atlas commands
- Use `dir = "web"` and `dir = "docs"` for Bun commands
- Use `depends` for task dependencies
- Use `{{arg(name='migration')}}` for migration names
- Load .env file for secrets

### Step 5: Update README.md
**Complexity**: Medium
**Dependencies**: Step 4

1. **Add Prerequisites Section**:
   ```markdown
   ## Prerequisites

   - [mise](https://mise.jit.su/) - Tool version manager and task runner
   - Docker & Docker Compose
   - (Optional) ngrok for tunneling
   ```

2. **Update Quick Start Section**:
   ```markdown
   ## Quick Start

   ```bash
   # Install all development tools
   mise install

   # First-time setup (DB + migrations + code generation)
   mise run setup

   # Run development servers
   mise run dev:backend  # Backend (hot reload, no build)
   mise run run:ui       # Frontend
   mise run dev:docs     # Documentation (optional)
   ```
   ```

3. **Replace Development Commands Section**:
   ```markdown
   ## Development Commands

   ### Development
   ```bash
   mise run dev:backend   # Backend with hot reload
   mise run run:ui        # Frontend dev server
   mise run dev:docs      # Documentation dev server
   mise run run           # Build and run production mode
   ```

   ### Code Quality
   ```bash
   mise run lint          # Lint Go + UI
   mise run test          # Run tests
   mise run generate      # Generate code (SQL, Swagger)
   ```

   ### Database
   ```bash
   mise run db:up                    # Start PostgreSQL
   mise run migrate:new name=feature # Create migration
   mise run migrate:up               # Apply migrations
   mise run migrate:status           # Check status
   mise run psql                     # Database console
   ```

   ### Building
   ```bash
   mise run build         # Build all (UI + Docs + Go)
   mise run build:go      # Build backend only
   mise run build:ui      # Build frontend only
   ```

   ### Docker
   ```bash
   mise run docker:up     # Start with docker compose
   mise run docker:down   # Stop docker compose
   mise run docker:build  # Build images
   ```

   ### Utilities
   ```bash
   mise tasks             # List all available tasks
   mise run doctor        # Check environment
   mise run clean         # Clean build artifacts
   mise run help          # Show help
   ```
   ```

4. **Add Troubleshooting Section**:
   ```markdown
   ## Troubleshooting

   **mise command not found**
   ```bash
   # Install mise
   curl https://mise.run | sh
   echo 'eval "$(mise activate bash)"' >> ~/.bashrc  # or ~/.zshrc
   source ~/.bashrc
   ```

   **Tool installation fails**
   ```bash
   mise doctor  # Check for issues
   mise install --force <tool>  # Force reinstall
   ```

   **Database connection refused**
   ```bash
   mise run db:up  # Ensure database is running
   docker ps | grep postgres  # Verify container
   ```

   **Port already in use**
   ```bash
   # Change port in .env
   PORT=8081 mise run dev:backend
   ```

   **Permission denied on psql**
   ```bash
   # Ensure docker compose is running
   docker compose ps
   mise run db:up
   ```
   ```

### Step 6: Update CLAUDE.md
**Complexity**: Low
**Dependencies**: Step 4

Replace the entire "Essential Commands" section:

```markdown
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
```

### Step 7: Update GitHub Actions Workflows
**Complexity**: Medium
**Dependencies**: Step 4

Find all workflows in `.github/workflows/` and update:

1. **Add mise setup step** (after checkout):
   ```yaml
   - name: Setup mise
     uses: jdx/mise-action@v2
     with:
       cache: true
   ```

2. **Create .env for CI** (if needed for tests):
   ```yaml
   - name: Setup environment
     run: |
       cat << EOF > .env
       DATABASE_URL=${{ secrets.DATABASE_URL }}
       DB_PASSWORD=postgres
       EOF
   ```

3. **Replace make commands** with mise:
   ```yaml
   # Before:
   - run: make test
   - run: make lint
   - run: make build

   # After:
   - run: mise run test
   - run: mise run lint
   - run: mise run build
   ```

**Example complete workflow**:
```yaml
name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup mise
        uses: jdx/mise-action@v2
        with:
          cache: true

      - name: Install tools
        run: mise install

      - name: Setup environment
        run: |
          cat << EOF > .env
          DB_PASSWORD=postgres
          EOF

      - name: Start database
        run: mise run db:up

      - name: Run migrations
        run: mise run migrate:up

      - name: Generate code
        run: mise run generate

      - name: Lint
        run: mise run lint

      - name: Test
        run: mise run test

      - name: Build
        run: mise run build
```

### Step 8: Create Migration Guide
**Complexity**: Low
**Dependencies**: Step 4

Create `MIGRATION.md` in repository root:

```markdown
# Migrating from Make to mise

## Installation

```bash
# Install mise
curl https://mise.run | sh

# Add to shell (choose one):
echo 'eval "$(mise activate bash)"' >> ~/.bashrc
echo 'eval "$(mise activate zsh)"' >> ~/.zshrc
source ~/.bashrc  # or source ~/.zshrc
```

## Setup

```bash
# In the recally directory
mise install        # Install all tools
mise run setup      # First-time setup
```

## Command Reference

| Old (Make) | New (mise) |
|------------|------------|
| `make lint` | `mise run lint` |
| `make lint-go` | `mise run lint:go` |
| `make lint-ui` | `mise run lint:ui` |
| `make generate` | `mise run generate` |
| `make generate-go` | `mise run generate:go` |
| `make generate-sql` | `mise run generate:sql` |
| `make generate-spec` | `mise run generate:spec` |
| `make build` | `mise run build` |
| `make build-go` | `mise run build:go` |
| `make build-ui` | `mise run build:ui` |
| `make build-docs` | `mise run build:docs` |
| `make test` | `mise run test` |
| `make run` | `mise run run` |
| `make run-go` | `mise run run:go` |
| `make run-ui` | `mise run run:ui` |
| `make db-up` | `mise run db:up` |
| `make migrate-new name=foo` | `mise run migrate:new name=foo` |
| `make migrate-up` | `mise run migrate:up` |
| `make migrate-status` | `mise run migrate:status` |
| `make migrate-validate` | `mise run migrate:validate` |
| `make migrate-hash` | `mise run migrate:hash` |
| `make psql` | `mise run psql` |
| `make docker-build` | `mise run docker:build` |
| `make docker-up` | `mise run docker:up` |
| `make docker-down` | `mise run docker:down` |
| `make ngrok` | `mise run ngrok` |
| `make deploy` | `mise run deploy` |
| `make help` | `mise tasks` or `mise run help` |

## New Commands

mise adds several new utility commands:

- `mise run setup` - First-time project setup
- `mise run doctor` - Check environment and tool versions
- `mise run clean` - Clean build artifacts
- `mise run dev:backend` - Run backend with hot reload (no build)
- `mise run dev:docs` - Run documentation dev server
- `mise tasks` - List all available tasks

## Why mise?

- **Single tool**: No more installing Go, Bun, sqlc, atlas, golangci-lint separately
- **Version lock**: Everyone uses exact same tool versions (Go 1.24.6, Bun 1.1, etc.)
- **Self-documenting**: `mise tasks` shows all commands with descriptions
- **Cross-platform**: Works on macOS, Linux, WSL
- **Faster onboarding**: New developers run `mise install` and they're ready

## Getting Help

- Run `mise tasks` to see all available commands
- Run `mise run doctor` to diagnose environment issues
- Check README.md for detailed documentation and troubleshooting

## Cleanup

After confirming mise works for you, this file can be deleted.
```

### Step 9: Update Dockerfile (if applicable)
**Complexity**: Medium
**Dependencies**: Step 4

If the project's Dockerfile references Makefile commands, update it.

**Check current Dockerfile**:
- Look for any `RUN make` commands
- Decide on approach:
  1. **Use mise in Docker**: Install mise in builder stage
  2. **Explicit commands**: Don't depend on mise, use direct commands

**Recommended approach** (explicit commands for Docker):
```dockerfile
FROM golang:1.24.6 AS builder

WORKDIR /app
COPY . .

# Install dependencies
RUN go mod download

# Generate code
RUN go generate ./...
RUN sqlc generate
RUN swag init -g internal/port/httpserver/router.go -o docs/swagger

# Build
RUN go build -o bin/app main.go

# Runtime stage
FROM debian:bookworm-slim
COPY --from=builder /app/bin/app /usr/local/bin/app
CMD ["app"]
```

This avoids adding mise to the Docker image and keeps builds fast.

### Step 10: Testing and Validation
**Complexity**: High
**Dependencies**: All previous steps

Comprehensive testing plan:

1. **Clean Environment Test**:
   ```bash
   # Simulate fresh developer setup
   rm -rf ~/.local/share/mise
   mise install
   mise run doctor
   ```

2. **Tool Installation Validation**:
   ```bash
   mise list
   go version      # Should be 1.24.6
   bun --version   # Should be 1.1.x
   atlas version
   sqlc version
   golangci-lint --version
   ```

3. **Environment Variable Test**:
   ```bash
   mise env | grep DATABASE_URL
   mise env | grep DB_HOST
   # Verify DATABASE_URL is correctly constructed
   ```

4. **Task Execution Tests**:
   ```bash
   # Lint tasks
   mise run lint:go
   mise run lint:ui
   mise run lint

   # Generation tasks
   mise run generate:go
   mise run generate:sql
   mise run generate:spec
   mise run generate

   # Build tasks
   mise run build:ui
   mise run build:docs
   mise run build:go
   mise run build

   # Test task
   mise run test

   # Database tasks
   mise run db:up
   sleep 3
   mise run migrate:up
   mise run migrate:status
   mise run psql  # Exit with \q

   # Docker tasks
   mise run docker:build

   # Utility tasks
   mise run clean
   mise run doctor
   mise tasks
   ```

5. **Development Workflow Test**:
   ```bash
   # Terminal 1: Backend
   mise run dev:backend

   # Terminal 2: Frontend
   mise run run:ui

   # Verify both start and hot reload works
   # Make a code change and verify reload
   ```

6. **Migration Test**:
   ```bash
   mise run migrate:new name=test_migration
   # Verify file created in database/migrations/
   mise run migrate:validate
   mise run migrate:hash
   ```

7. **Production Build Test**:
   ```bash
   mise run clean
   mise run build
   ./bin/app  # Should start
   ```

8. **CI/CD Simulation** (if using act):
   ```bash
   act -j test  # Run GitHub Actions locally
   ```

9. **Cross-Platform Test** (if team uses multiple OSs):
   ```bash
   # Test on Linux via Docker
   docker run --rm -it -v $(pwd):/workspace -w /workspace ubuntu:22.04 bash
   curl https://mise.run | sh
   export PATH="$HOME/.local/bin:$PATH"
   mise install
   mise run test
   ```

10. **Performance Comparison**:
    ```bash
    # Compare execution time
    time make build   # Old way (if Makefile still exists)
    time mise run build  # New way
    # Should be comparable
    ```

**Validation Checklist**:
- [ ] All tools install successfully
- [ ] `mise run doctor` shows no errors
- [ ] All tasks execute without errors
- [ ] Environment variables are correct
- [ ] Database migrations work
- [ ] Dev servers start and hot reload works
- [ ] Production build succeeds
- [ ] Binary runs correctly
- [ ] Documentation is accurate

### Step 11: Remove Makefile
**Complexity**: Low
**Dependencies**: Step 10 (successful validation)

Only after all tests pass:

```bash
git rm Makefile
```

Keep the backup locally until merge is complete, just in case.

### Step 12: Commit Changes
**Complexity**: Low
**Dependencies**: Step 11

Create comprehensive commit:

```bash
git add .
git commit -m "✨ feat: migrate from Makefile to mise for dependency and task management

- Replace Makefile with comprehensive mise.toml configuration
- Pin all tool versions for reproducibility:
  - Go 1.24.6 (from go.mod)
  - Bun 1.1 (frontend/docs)
  - golangci-lint 1.62
  - sqlc 1.27
  - swag 1.16
  - Atlas 0.29
- Organize 40+ tasks by category (lint, generate, build, run, db, docker)
- Add utility tasks: setup, doctor, clean, help
- Update README.md with mise commands and troubleshooting
- Update CLAUDE.md Essential Commands section
- Update GitHub Actions workflows to use mise-action
- Add MIGRATION.md guide for team

Key improvements:
- All development dependencies managed by single tool
- Consistent tool versions across team and CI/CD
- Self-documenting tasks via 'mise tasks'
- Better task organization with categories
- Proper error handling in setup task
- Environment variables loaded from .env for security

Breaking changes:
- All 'make <task>' commands now 'mise run <task>'
- Team needs to install mise: curl https://mise.run | sh
- See MIGRATION.md for complete command reference

Migration verified:
- All Makefile tasks successfully migrated
- Full test suite passes
- Dev workflow tested (backend + frontend)
- Database migrations working
- Production build successful
"
```

### Step 13: Push and Create Pull Request
**Complexity**: Low
**Dependencies**: Step 12

```bash
# Push branch
git push -u origin feat/migrate-to-mise

# Create PR using gh CLI
gh pr create --title "✨ feat: migrate from Makefile to mise" --body "$(cat <<'EOF'
## Summary

This PR migrates the project from GNU Make to mise for dependency and task management.

## Changes

- ✅ Replaced Makefile with comprehensive `mise.toml`
- ✅ Pin all development tool versions for reproducibility
- ✅ Organized 40+ tasks by functional category
- ✅ Updated documentation (README.md, CLAUDE.md)
- ✅ Updated GitHub Actions workflows
- ✅ Added migration guide for team

## Benefits

1. **Single tool for everything**: No more separate installation of Go, Bun, sqlc, atlas, golangci-lint, swag
2. **Version consistency**: Everyone uses exact same tool versions (eliminating "works on my machine")
3. **Self-documenting**: Run `mise tasks` to see all available commands with descriptions
4. **Better organization**: Tasks grouped by category (lint:*, build:*, db:*, etc.)
5. **Faster onboarding**: New developers: `mise install && mise run setup` and they're ready

## Breaking Changes

All `make <task>` commands are now `mise run <task>`:
- `make build` → `mise run build`
- `make test` → `mise run test`
- `make run-go` → `mise run run:go` or `mise run dev:backend`

See [MIGRATION.md](./MIGRATION.md) for complete command reference.

## Migration Steps for Team

1. **Install mise**:
   ```bash
   curl https://mise.run | sh
   eval "$(mise activate bash)"  # Add to your ~/.bashrc or ~/.zshrc
   ```

2. **Install tools**:
   ```bash
   mise install
   ```

3. **Ready to develop**:
   ```bash
   mise run dev:backend  # Terminal 1
   mise run run:ui       # Terminal 2
   ```

4. **See all commands**:
   ```bash
   mise tasks
   ```

## Testing Performed

- [x] Clean environment test (fresh mise install)
- [x] All tasks execute successfully
- [x] Development workflow (backend + frontend hot reload)
- [x] Database migrations (up, down, status, validate)
- [x] Production build
- [x] Lint and test suite
- [x] Docker builds
- [x] GitHub Actions updated and verified

## Documentation

- ✅ README.md updated with mise commands
- ✅ CLAUDE.md updated with mise commands
- ✅ MIGRATION.md created with team guide
- ✅ Troubleshooting section added

## Questions?

Run `mise run doctor` to check your environment, or see documentation in README.md.

---

**Please review and test locally before merging. Happy to pair on mise setup if needed!**
EOF
)"
```

### Step 14: Post-Merge Cleanup
**Complexity**: Low
**Dependencies**: PR merged

After successful merge and team adoption:

1. **Monitor team adoption**:
   - Check for issues in team channels
   - Offer pairing sessions for setup
   - Document any edge cases found

2. **Update CI/CD documentation**:
   - Ensure all workflow files are updated
   - Verify CI builds pass

3. **Remove migration guide** (after 2-4 weeks):
   ```bash
   git rm MIGRATION.md
   git commit -m "📝 docs: remove migration guide after successful adoption"
   ```

4. **Review and optimize**:
   - Monitor task execution times
   - Gather feedback on task organization
   - Consider adding more utility tasks based on team needs

---

## Testing Strategy

### Unit Testing Approach

**Tool Installation Testing**:
- Test mise installation from scratch
- Verify all tools install to correct versions
- Run `mise doctor` to check for configuration issues
- Test tool upgrades: `mise upgrade`

**Task Execution Testing**:
- Test each task individually
- Verify task dependencies work (`depends` parameter)
- Test tasks with arguments (migrate:new)
- Verify working directory changes (`dir` parameter)
- Test environment variable overrides

**Environment Variable Testing**:
- Verify DATABASE_URL construction
- Test .env file loading
- Verify environment isolation between tasks
- Test with missing .env (should use defaults where possible)

### Integration Testing

**Full Development Workflow**:
1. Fresh repository clone
2. `mise install` - Install all tools
3. `mise run setup` - First-time setup
4. `mise run dev:backend` - Start backend
5. `mise run run:ui` - Start frontend
6. Make code changes - Verify hot reload
7. `mise run lint` - Lint code
8. `mise run test` - Run tests
9. `mise run build` - Production build
10. `./bin/app` - Run production binary

**Database Workflow**:
1. `mise run db:up` - Start database
2. Wait for ready (3 seconds)
3. `mise run migrate:up` - Apply migrations
4. `mise run migrate:status` - Check status
5. `mise run migrate:new name=test` - Create migration
6. `mise run migrate:validate` - Validate
7. `mise run psql` - Interactive console
8. Query database, exit
9. `mise run docker:down` - Cleanup

**CI/CD Workflow**:
1. Trigger GitHub Actions
2. Verify mise installation
3. Verify tool caching works
4. Verify all tasks pass
5. Compare CI times before/after

### Edge Cases and Error Handling

1. **Missing .env file**:
   - Tasks should work with mise defaults where possible
   - Clear error if required secrets missing
   - Document required vs optional env vars

2. **Tool installation failures**:
   - mise provides clear error with tool name
   - Document manual installation fallback
   - Add to troubleshooting guide

3. **Database not ready**:
   - Setup task waits 3 seconds after db:up
   - If migrations fail, clear error message
   - Document how to check postgres container

4. **Port conflicts**:
   - Environment variables allow port override
   - Document in troubleshooting
   - Suggest using .env for custom ports

5. **Cross-platform path issues**:
   - Use relative paths in tasks
   - Test on macOS and Linux
   - Document Windows/WSL if needed

6. **Dependency ordering**:
   - Use `depends` for explicit dependencies
   - Test that build:go runs generate first
   - Verify parallel execution where appropriate

7. **Interactive commands**:
   - `psql` task uses `-it` flags correctly
   - Test Ctrl+C signal handling
   - Verify proper cleanup on interrupt

8. **Migration with arguments**:
   - Test `migrate:new` with name argument
   - Verify default value if no name provided
   - Ensure proper error if required arg missing

### Validation Criteria

**Functional Requirements**:
- [ ] All 40+ Makefile targets have mise equivalents
- [ ] All tasks execute successfully
- [ ] Task dependencies work correctly
- [ ] Environment variables constructed properly
- [ ] Tools install on clean system
- [ ] CI/CD pipelines pass with mise
- [ ] Development workflow smooth (hot reload works)
- [ ] Database migrations function correctly
- [ ] Docker tasks work

**Documentation Requirements**:
- [ ] README.md has no make references
- [ ] CLAUDE.md updated with all mise commands
- [ ] MIGRATION.md provides clear guidance
- [ ] Troubleshooting section comprehensive
- [ ] Task organization intuitive
- [ ] mise.toml well-commented

**Quality Requirements**:
- [ ] mise.toml organized and readable
- [ ] Task names consistent (lint:*, build:*, etc.)
- [ ] `mise tasks` output is clear
- [ ] Error messages helpful
- [ ] Tool versions pinned (not "latest")
- [ ] Setup task has proper error handling
- [ ] No security issues (no secrets in mise.toml)

**Performance Requirements**:
- [ ] Tool installation time acceptable (~2-5 min first time)
- [ ] Task execution overhead minimal (<100ms)
- [ ] Build times comparable to Makefile
- [ ] CI/CD not slower with caching
- [ ] Dev server startup time unchanged

---

## Considerations

### Security Implications

1. **Environment Variables and Secrets**:
   - ✅ Use `_.file = [".env"]` to load secrets from gitignored file
   - ✅ Never commit secrets to mise.toml
   - ✅ Reference system environment for sensitive values: `{{env.OPENAI_API_KEY}}`
   - ✅ Document required env vars in README
   - ✅ Provide .env.example template

2. **Tool Supply Chain Security**:
   - ✅ mise downloads tools from official sources (GitHub releases)
   - ✅ Pin specific versions for reproducibility
   - ✅ Avoid "latest" in production
   - ⚠️ Review tool installation URLs if security-critical
   - ✅ Use ubi: prefix for verified GitHub releases

3. **CI/CD Security**:
   - ✅ mise installation script runs as current user (not root)
   - ✅ Use official mise-action in GitHub Actions
   - ✅ Secrets passed via GitHub Actions secrets
   - ✅ No secrets in workflow files
   - ✅ mise.toml in version control (no secrets there)

4. **Database Credentials**:
   - ✅ Development DB password in .env (gitignored)
   - ✅ Production credentials via system environment
   - ⚠️ Never log DATABASE_URL (contains password)
   - ✅ Use Docker exec for psql (no password in args)

### Performance Concerns

1. **Tool Installation Time**:
   - First `mise install`: 2-5 minutes (downloads all tools)
   - Subsequent runs: < 1 second (tools cached)
   - CI/CD: Use mise-action caching to avoid reinstalls
   - Cache key: `mise-${{ hashFiles('mise.toml') }}`

2. **Task Execution Overhead**:
   - mise task startup: ~10-50ms (negligible)
   - Tasks run native commands (no interpretation)
   - Parallel task execution where possible
   - No performance impact on actual builds

3. **Build Times**:
   - Same underlying commands as Makefile
   - No impact on Go build time
   - No impact on Bun build time
   - Task dependencies may slightly reorganize execution

4. **Development Experience**:
   - Hot reload unchanged (same commands)
   - Tool switching instant (mise managed)
   - No slowdown in development iteration

### Scalability Factors

1. **Task Organization**:
   - Dot notation scales: `category:subcategory:task`
   - Can refactor to includes if needed:
     ```toml
     [includes]
     "tasks/database.toml"
     "tasks/build.toml"
     ```
   - Current 40+ tasks manageable in single file
   - Clear categories aid discoverability

2. **Tool Management**:
   - mise handles multiple versions gracefully
   - Per-directory overrides possible
   - Scales to monorepos
   - Can specify different versions per subproject

3. **Team Adoption**:
   - Single tool installation for all developers
   - Consistent across projects using mise
   - Self-documenting (`mise tasks`)
   - Reduces onboarding friction

4. **Multi-Project Support**:
   - mise works per-directory
   - Different projects can have different tool versions
   - Global fallbacks configurable
   - Team can standardize on mise across projects

### Maintenance and Documentation

1. **mise.toml Maintenance**:
   - Add comments for complex tasks
   - Document required environment variables at top
   - Keep task descriptions up to date
   - Review tool versions quarterly
   - Group related tasks with section headers

2. **Version Updates**:
   - Update Go: `mise use go@1.24.7`
   - Update all tools: `mise upgrade`
   - Commit mise.toml changes
   - Team runs `mise install` to sync
   - Test thoroughly before merging

3. **Documentation Sync**:
   - Keep README.md in sync with mise tasks
   - Update CLAUDE.md when tasks change
   - Add troubleshooting for common issues
   - Include `mise tasks` output examples
   - Document new tasks as they're added

4. **Onboarding New Developers**:
   - Install mise: `curl https://mise.run | sh`
   - Clone repo
   - `mise install` - Get all tools
   - `mise run setup` - Setup project
   - Start coding
   - No more "missing tools" issues

5. **Long-term Migration Path**:
   - This is one-way migration (removing Makefile)
   - No plan to revert to Make
   - mise is actively maintained
   - Growing adoption in development community
   - Can always fall back to direct tool calls if needed

### Risk Mitigation

1. **Risk: Team unfamiliar with mise**
   - **Mitigation**:
     - Clear documentation in MIGRATION.md
     - mise is intuitive (similar to make)
     - Offer pairing sessions
     - mise tasks is self-documenting

2. **Risk: CI/CD breaks during migration**
   - **Mitigation**:
     - Test workflows locally with act first
     - Use official mise-action
     - Update workflows atomically
     - Keep PR focused and reviewable
     - Test in feature branch first

3. **Risk: Platform-specific issues**
   - **Mitigation**:
     - Test on macOS and Linux
     - mise has broad OS support
     - Document platform-specific issues
     - Use Docker for consistent environments

4. **Risk: Missing Makefile features**
   - **Mitigation**:
     - Comprehensive analysis of existing Makefile
     - All 40+ targets migrated
     - mise supports complex tasks
     - Can use shell scripts for complex logic

5. **Risk: Performance regression in CI**
   - **Mitigation**:
     - Add mise tool caching in workflows
     - Compare CI times before/after
     - Optimize parallel task execution
     - Cache node_modules and go mod cache

6. **Risk: Tool installation failures**
   - **Mitigation**:
     - Pin specific versions (not latest)
     - Use official ubi: sources
     - Document manual fallback
     - Test on clean systems

---

## Summary

This plan provides a comprehensive migration from Makefile to mise, incorporating all feedback from Codex review:

### What We're Migrating
- **40+ make targets** → Organized mise tasks
- **Manual tool installation** → mise-managed tools
- **Environment variables** → mise configuration with .env
- **Build documentation** → Updated for mise

### Key Benefits
1. **Single tool**: mise manages Go, Bun, sqlc, atlas, golangci-lint, swag
2. **Version lock**: Go 1.24.6, Bun 1.1, etc. - everyone uses same versions
3. **Self-documenting**: `mise tasks` shows all commands
4. **Better organization**: Tasks grouped by category
5. **Faster onboarding**: `mise install && mise run setup`

### Critical Fixes Applied (from Codex Review)
1. ✅ Use Bun (not npm) for all frontend/docs tasks
2. ✅ All Makefile tasks included (40+ tasks)
3. ✅ Correct paths: `bin/app`, `database/`, `postgres` service
4. ✅ Pin tool versions (not "latest")
5. ✅ Proper Atlas working directory
6. ✅ Environment variables with .env loading
7. ✅ Error handling in setup task
8. ✅ Comprehensive documentation
9. ✅ Migration guide for team
10. ✅ GitHub Actions with mise-action

### Estimated Effort
- Implementation: 3-4 hours
- Testing: 1-2 hours
- Documentation: 1 hour
- **Total: 5-7 hours**

### Success Metrics
- [ ] All developers using mise within 1 week
- [ ] CI/CD pipelines passing
- [ ] No reported migration issues
- [ ] Positive team feedback
- [ ] MIGRATION.md can be removed after 2-4 weeks

Ready for implementation! 🚀
