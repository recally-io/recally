# Recally ✨
**Your AI-Powered Memory Assistant for Digital Content**

![logo](./web/public/logo.svg)

Never lose track of valuable content again. Recally helps you capture, organize, and rediscover knowledge with AI-powered efficiency.

## ⚡ Key Features

### 📚 Content Capture & Types
- **🌐 Web Content**
  - ✅ Browser extension for one-click saving
  - ✅ Telegram bot integration for mobile capture
  - ✅ Direct in-app article addition
  - ✅ Automatic markdown conversion
  - ✅ Smart CORS image proxy

- **💬 Social Media**
  - ✅ Twitter thread unrolling and saving
  - ⌛ Instagram post archiving
  - ⌛ LinkedIn article saving

- **📱 Rich Media**
  - ⌛ YouTube video saving & transcription
  - ⌛ PDF document storage & analysis
  - ⌛ Podcast archiving & transcription
  - ⌛ EPUB book organization

### 🧠 AI-Powered Knowledge Management
- **🤖 Automated Processing**
  - ✅ Smart summarization with key points
  - ✅ Intelligent tag suggestions
  - ✅ Content categorization
  - ⌛ RAG-powered knowledge base
  - ⌛ Interactive document Q&A

- **🔍 Discovery & Search**
  - ✅ Lightning-fast full-text search
  - ✅ Smart filtering by tags/domains
  - ⌛ Semantic search across content
  - ⌛ Related content suggestions

### 🔌 Integration Ecosystem
- **🌍 Browser Extensions**
  - ✅ [![Chrome](https://img.shields.io/badge/Chrome-Extension-brightgreen?logo=googlechrome)](https://chrome.google.com/webstore/detail/heblpkdddipfjdpdgikoledoecohoepp)
  - ✅ [![Firefox](https://img.shields.io/badge/Firefox-Add_on-FF7139?logo=firefoxbrowser)](https://addons.mozilla.org/addon/recally-clipper/)
  - ⌛ Safari extension

- **📱 Mobile & Messaging**
  - ✅ [![Telegram Bot](https://img.shields.io/badge/Telegram-RecallyReaderBot-2CA5E0?logo=telegram)](https://t.me/RecallyReaderBot)
  - ⌛ Mobile apps (iOS/Android)

- **📝 Note-Taking Apps**
  - ⌛ Notion sync
  - ⌛ Obsidian plugin

- **📰 Content Sources**
  - ⌛ RSS feed integration
  - ⌛ Newsletter management
  - ⌛ Email forwarding

## 🚀 Get Started in 60 Seconds

### Cloud Version (Recommended)
👉 [recally.io](https://recally.io)  
- Instant access with Google/GitHub login
- Free during beta (no credit card required)
- Always up-to-date

### Self-Hosted Option
For full control over your data:
```bash
git clone https://github.com/recally-io/recally
cd recally
cp env.example .env  # Set OpenAI key & DB credentials
docker compose up -d
```
Access at `http://localhost:1323`

> **Note:** Requires [Docker](https://docs.docker.com/get-docker/) and [OpenAI API key](https://platform.openai.com/api-keys)

## 🛠 Developer Zone

### Prerequisites

- [mise](https://mise.jit.su/) - Tool version manager and task runner
- Docker & Docker Compose
- (Optional) ngrok for tunneling

### Tech Stack
**Backend**
[![Go](https://img.shields.io/badge/Go-1.24.6-00ADD8?logo=go)](https://go.dev/)
[![Echo](https://img.shields.io/badge/Echo-v4.11-blue)](https://echo.labstack.com/)
[![ParadeDB](https://img.shields.io/badge/ParadeDB-1.0-orange)](https://www.paradedb.com/)

**Frontend**
[![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)](https://react.dev/)
[![Vite](https://img.shields.io/badge/Vite-5.0-646CFF?logo=vite)](https://vitejs.dev/)

**AI**
[![OpenAI](https://img.shields.io/badge/OpenAI-GPT4-purple?logo=openai)](https://openai.com/)
[![LocalAI](https://img.shields.io/badge/Option-Ollama-blue)](https://ollama.com/)

### Quick Start for Contributors

1. **Install mise**:
   ```bash
   curl https://mise.run | sh
   eval "$(mise activate bash)"  # Add to your ~/.bashrc or ~/.zshrc
   ```

2. **Clone and setup**:
   ```bash
   git clone https://github.com/recally-io/recally
   cd recally
   cp env.example .env  # Configure environment variables

   # Install all development tools
   mise install

   # First-time setup (DB + migrations + code generation)
   mise run setup
   ```

3. **Start developing**:
   ```bash
   # Terminal 1: Backend (hot reload)
   mise run dev:backend

   # Terminal 2: Frontend
   mise run run:ui
   ```

4. Check our [Good First Issues](https://github.com/recally-io/recally/contribute)

### Development Commands

#### Development
```bash
mise run dev:backend   # Backend with hot reload
mise run run:ui        # Frontend dev server
mise run dev:docs      # Documentation dev server
mise run run           # Build and run production mode
```

#### Code Quality
```bash
mise run lint          # Lint Go + UI
mise run test          # Run tests
mise run generate      # Generate code (SQL, Swagger)
```

#### Database
```bash
mise run db:up                    # Start PostgreSQL
mise run migrate:new name=feature # Create migration
mise run migrate:up               # Apply migrations
mise run migrate:status           # Check status
mise run psql                     # Database console
```

#### Building
```bash
mise run build         # Build all (UI + Docs + Go)
mise run build:go      # Build backend only
mise run build:ui      # Build frontend only
```

#### Docker
```bash
mise run docker:up     # Start with docker compose
mise run docker:down   # Stop docker compose
mise run docker:build  # Build images
```

#### Utilities
```bash
mise tasks             # List all available tasks
mise run doctor        # Check environment
mise run clean         # Clean build artifacts
mise run help          # Show help
```

### Troubleshooting

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

## 📜 Documentation
Explore our comprehensive guides:
- [Documentation](https://recally.io/docs/)
- [API Reference](https://recally.io/swagger/index.html)

## Comparison with Similar Tools

| Feature | Recally | Pocket | Instapaper | Shiori | Omnivore | Readwise Reader | Hoarder |
|---------|---------|--------|------------|---------|-----------|----------------|---------|
| 🌐 Web Article Saving | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| 🧵 Twitter Thread Support | ✅ | ❌ | ❌ | ❌ | ✅ | ✅ | ❌ |
| 🤖 AI Summarization | ✅ | ❌ | ❌ | ❌ | ✅ | ✅ | ✅ |
| 🏷️ Smart Tagging | ✅ | ❌ | ❌ | ❌ | ✅ | ✅ | ✅ |
| 🖼️ CORS Image Proxy | ✅ | ❌ | ❌ | ✅ | ✅ | ✅ | ❌ |
| 🛠️ Self-Hosted Option | ✅ | ❌ | ❌ | ✅ | ✅ | ❌ | ✅ |
| 📱 Mobile Support | ✅  | ✅ | ✅ | ❌ | ✅ | ✅ |  ✅ |

## License
- **Non-commercial**: [AGPLv3](LICENSE)
- **Commercial**: Contact [sales@recally.io](mailto:sales@recally.io) for enterprise licensing

---

> Made with ♥ by Recally Team  
> Proudly open-core since 2024
