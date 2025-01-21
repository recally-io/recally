# Recally ✨
**Your AI-Powered Memory Assistant for Digital Content**

![logo](./web/public/logo.svg)

Never lose track of valuable content again. Recally helps you capture, organize, and rediscover knowledge with AI-powered efficiency.

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

## 🔥 Why Recally?

### Core Features
| Category | Features |
|----------|----------|
| 📥 Capture | One-click web saves • [Browser extensions](https://github.com/recally-io/recally-clipper) • [Telegram bot](https://t.me/RecallyReaderBot) • PDF import (soon) |
| 🧠 Intelligence | AI summarization • Smart tagging • Semantic search • Document Q&A (soon) |
| 🛡 Privacy | Self-hostable • Zero tracking • Open-source core |

### Unique Advantages
- **AI That Understands Context**  
  GPT-4 powered analysis that goes beyond keyword matching
- **Multi-Source Support**  
  Articles, YouTube videos, podcasts, PDFs - all in one place
- **True Ownership**  
  Export all data anytime • No lock-in or ads

## 📱 Capture Content Anywhere

### Browser Extensions
[![Chrome](https://img.shields.io/badge/Chrome-Extension-brightgreen?logo=googlechrome)](https://chrome.google.com/webstore/detail/heblpkdddipfjdpdgikoledoecohoepp)
[![Firefox](https://img.shields.io/badge/Firefox-Add_on-FF7139?logo=firefoxbrowser)](https://addons.mozilla.org/addon/recally-clipper/)

Features:
- Save pages with original formatting
- Highlight key sections

### Telegram Bot
[![Telegram Bot](https://img.shields.io/badge/Telegram-RecallyReaderBot-2CA5E0?logo=telegram)](https://t.me/RecallyReaderBot)

Send any link to:
- Save instantly to your library
- Get 3-sentence AI summary

## 🛠 Developer Zone

### Tech Stack
**Backend**  
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Echo](https://img.shields.io/badge/Echo-v4.11-blue)](https://echo.labstack.com/)
[![ParadeDB](https://img.shields.io/badge/ParadeDB-1.0-orange)](https://www.paradedb.com/)

**Frontend**  
[![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)](https://react.dev/)
[![Vite](https://img.shields.io/badge/Vite-5.0-646CFF?logo=vite)](https://vitejs.dev/)

**AI**  
[![OpenAI](https://img.shields.io/badge/OpenAI-GPT4-purple?logo=openai)](https://openai.com/)
[![LocalAI](https://img.shields.io/badge/Option-Ollama-blue)](https://ollama.com/)

### Contribution Guide
1. Fork & clone repo
2. Set up dev environment:
```bash
make run  # Starts both backend and frontend with hot-reload
```
3. Check our [Good First Issues](https://github.com/recally-io/recally/contribute)

## 📜 Documentation
Explore our comprehensive guides:
- [Documentation](https://recally.io/docs/)
- [API Reference](https://recally.io/swagger/index.html)

## Similar Tools

While we love Recally, here are some other great options:
- [Shiori](https://github.com/go-shiori/shiori) - Great for CLI lovers
- [Omnivore](https://omnivore.app) - RIP, you were awesome
- [Pocket](https://getpocket.com) - The OG save-for-later app
- [Readwise Reader](https://readwise.io) - The king of highlights
- [Hoarder](https://github.com/hoarder-app/hoarder) - Self-hostable bookmark manager with AI features
- [Instapaper](https://www.instapaper.com) - Clean, minimalist read-it-later service

## License
- **Non-commercial**: [AGPLv3](LICENSE)
- **Commercial**: Contact [sales@recally.io](mailto:sales@recally.io) for enterprise licensing

---

> Made with ♥ by Recally Team
> Proudly open-core since 2024
