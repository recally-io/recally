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
- `mise run migrate:down` - Rollback last migration
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
