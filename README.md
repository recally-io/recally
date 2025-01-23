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
