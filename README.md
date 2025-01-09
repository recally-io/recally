# Recally

[Recally.io](https://recally.io) is your personal knowledge engine, designed to help you collect, organize, and remember the digital content that matters most. Save articles, videos, and podcasts with a single click, and let Recally’s smart tools help you recall and connect ideas when you need them. Say goodbye to information overload and hello to effortless learning.

![logo](./web/public/logo.svg)

## 🚀 Key Features

### 📥 Content Integration
- **Universal Content Support**
  - Save web articles with one click
  - [ ] Import PDFs, EPUBs, and documents
  - [ ] RSS feed reader
  - [ ] YouTube and Bilibili video synchronization
  - [ ] Podcast episode archiving
  - [ ] Newsletter integration

### 🤖 AI-Powered Intelligence
- **Smart Content Processing**
  - Automatic summarization and key point extraction
  - Intelligent tagging and topic detection
  - Custom AI prompts for personalized analysis
  - [ ] Multi-language translation support
  - [ ] Interactive document chat and Q&A

### 🔍 Advanced Search & Discovery
- **Powerful Search Capabilities**
  - Real-time full-text search
  - Filter by source, tags, dates, and more
  - [ ] Semantic similarity search

### 📚 Content Management
- **Robust Content Processing**
  - Multiple fetcher support (HTTP, Jina, Headless browser)
  - Automatic image preservation
  - Smart content cleanup and formatting

### 🔐 Privacy & Security
- **User-First Design**
  - Self-hosted option available
  - No third-party tracking
  - Complete data ownership
  - [ ] Security audits
  - [ ] End-to-end encryption for sensitive data

## 🛠 Installation

```bash

# Clone the repository
git clone https://github.com/vaayne/recally

# Change directory
cd recally

# edit the .env file as needed
cp env.example .env
# vim .env

# Run the application
docker compose up -d

# Access the application
open http://localhost:1323
```

## 🏗 Tech Stack

### Backend
- **API Server**: [Echo](https://github.com/labstack/echo) - High performance, minimalist Go web framework
- **Job Queue**: [River](https://github.com/riverqueue/river) - Background job processing
- **Database**: 
  - [ParadeDB](https://github.com/paradedb/paradedb) - Postgres with Search and Analytics 
- **Tools**:
  - [Migrate](https://github.com/golang-migrate/migrate) - Database migrations
  - [Sqlc](https://github.com/sqlc-dev/sqlc) - Type-safe SQL

### Frontend
- **Framework**: [React](https://github.com/facebook/react) - UI development
- **Build Tool**: [Vite](https://github.com/vitejs/vite) - Next generation frontend tooling
- **Styling**: 
  - [TailwindCSS](https://github.com/tailwindlabs/tailwindcss) - Utility-first CSS
  - [shadcn/ui](https://github.com/shadcn-ui/ui) - Accessible components

### AI Integration
- [OpenAI](https://openai.com/) and compatible models for:
  - Content analysis
  - Tag generation
  - Semantic search

## 🤝 Contributing

We welcome contributions!

## 🙏 Thanks

This project stands on the shoulders of giants. Special thanks to these amazing open-source projects:

- [go-readability](https://github.com/go-shiori/go-readability) - Clean article extraction
- [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) - HTML to Markdown conversion

And many other open-source projects that make this possible! 💖

## Alternatives

- [Shiori](https://github.com/go-shiori/shiori) - Simple, CLI-focused bookmark manager written in Go
- [Hoarder](https://github.com/hoarder-app/hoarder) - Self-hostable bookmark manager with AI features
- [Omnivore](https://omnivore.app) - Open-source read-it-later app with social features (closed now)
- [Pocket](https://getpocket.com) - Popular commercial bookmarking service by Mozilla
- [Instapaper](https://www.instapaper.com) - Clean, minimalist read-it-later service

## 📝 License

See the [LICENSE](LICENSE) file for details.

- **Non-commercial Use**: Free under MIT License terms
- **Commercial Use**: Contact [support@recally.io](mailto:support@recally.io) for permission

