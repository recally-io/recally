# Recally

Ever felt overwhelmed trying to keep track of all the interesting stuff you find online? That's exactly why we built Recally. It's a simple tool that helps you save and recall the content that matters to you, powered by AI to make it actually useful.

![logo](./web/public/logo.svg)

## 🚀 Quick Start

1. **Using the Cloud Version**
   - Visit [recally.io](https://recally.io) to get started immediately
   - Free during beta period
   - No installation required

2. **Self-Hosting**
   ```bash
   git clone https://github.com/recally-io/recally
   cd recally
   cp env.example .env    # Configure your settings
   docker compose up -d
   ```
   Visit http://localhost:1323 to start using your instance.

## What's Special About Recally?

### 🎯 Save Anything, Find Everything
- One-click saving of articles and web pages (and they actually look good when saved!)
- Coming soon: PDF imports, YouTube videos, and podcast episodes
- Smart search that actually understands what you're looking for

### 🤖 AI That Makes Sense
We're not just throwing AI in because it's trendy. Recally uses AI to:
- Create quick summaries so you remember why you saved something
- Suggest tags that actually make sense
- Help you connect ideas across your saved content
- Let you chat with your documents (coming soon!)

### 🔒 Your Content, Your Control
- Self-host if you want to (yes, we actually made this easy)
- No sneaky tracking or data sharing
- Keep everything organized your way

## Try It Out

Getting started with Recally is super easy:

### 🤖 Quick Save with Telegram
Just start chatting with our [RecallyReader](https://t.me/RecallyReaderBot) Telegram bot:
- Send any link to save articles and web pages
- Get instant AI-powered summaries
- Access your saved content anywhere

### 🌐 Web Experience
Head over to [recally.io](https://recally.io) to unlock the full potential:
- Beautiful reading interface
- Smart organization with AI-suggested tags
- Advanced search capabilities
- Free during beta, with new features added regularly

### 🔗 Browser Extensions
Save content with just one click using our browser extensions [Recally Clipper](https://github.com/recally-io/recally-clipper):

- [Chrome Extension](https://chrome.google.com/webstore/detail/heblpkdddipfjdpdgikoledoecohoepp)
- [Firefox Add-on](https://addons.mozilla.org/addon/recally-clipper/)

## 🛠 Development Setup

1. **Prerequisites**
   - Go 1.21+
   - Node.js 18+
   - Docker and Docker Compose
   - OpenAI API key (or compatible model)

2. **Local Development**
   ```bash
   # Backend
   cd api
   go mod download
   go run main.go

   # Frontend
   cd web
   npm install
   npm run dev
   ```

3. **Environment Variables**
   Key configurations in `.env`:
   ```env
   OPENAI_API_KEY=your_key_here
   DB_URL=postgresql://user:pass@localhost:5432/recally
   ```

## 📚 Documentation

- https://recally.io/docs/

Our REST API is documented using OpenAPI/Swagger:
- Local: http://localhost:1323/swagger/index.html
- Cloud: https://recally.io/swagger/index.html

## Under the Hood

We've chosen our tech stack carefully to make Recally fast, reliable, and easy to maintain:

### Backend
- [Echo](https://github.com/labstack/echo) for the API (because Go is fast and Echo is simple)
- [River](https://github.com/riverqueue/river) for job processing (rock-solid queue management)
- [ParadeDB](https://github.com/paradedb/paradedb) (Postgres + search that actually works)

### Frontend
- [React](https://github.com/facebook/react) (you know it, you love it)
- [Vite](https://github.com/vitejs/vite) (because waiting for builds is no fun)
- [TailwindCSS](https://github.com/tailwindlabs/tailwindcss) + [shadcn/ui](https://github.com/shadcn-ui/ui) (beautiful UI without the bloat)

### AI Magic
We use OpenAI (and compatible models) to make your content more useful through:
- Smart summaries
- Intelligent tagging
- Semantic search that understands context

## Want to Help?

We love contributions! Whether it's:
- Finding bugs
- Suggesting features
- Improving docs
- Adding code

Just jump in! Check our [issues](https://github.com/recally-io/recally/issues) or start a [PR](https://github.com/recally-io/recally/pulls). We're friendly, promise!

## Standing on Giants

Huge thanks to these amazing projects that make Recally possible:
- [go-readability](https://github.com/go-shiori/go-readability) for making saved articles beautiful
- [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) for clean content conversion
- And many others that deserve a beer 🍺

## Similar Tools

While we love Recally, here are some other great options:
- [Shiori](https://github.com/go-shiori/shiori) - Great for CLI lovers
- [Omnivore](https://omnivore.app) - RIP, you were awesome
- [Pocket](https://getpocket.com) - The OG save-for-later app
- [Readwise Reader](https://readwise.io) - The king of highlights
- [Hoarder](https://github.com/hoarder-app/hoarder) - Self-hostable bookmark manager with AI features
- [Instapaper](https://www.instapaper.com) - Clean, minimalist read-it-later service

## 📊 Status

[![Go Report Card](https://goreportcard.com/badge/github.com/recally-io/recally)](https://goreportcard.com/report/github.com/recally-io/recally)
[![License](https://img.shields.io/badge/license-AGPL%20v3-blue.svg)](LICENSE)

## License

Free for personal use under [GNU AGPLv3](https://choosealicense.com/licenses/agpl-3.0/). For commercial stuff, drop us a line at [support@recally.io](mailto:support@recally.io).
