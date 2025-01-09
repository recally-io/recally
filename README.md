# Recally

Ever felt overwhelmed trying to keep track of all the interesting stuff you find online? That's exactly why we built Recally. It's a simple tool that helps you save and recall the content that matters to you, powered by AI to make it actually useful.

![logo](./web/public/logo.svg)

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

Head over to [recally.io](https://recally.io) to give it a spin. It's free while in beta, and we're adding new stuff almost daily based on user feedback.

## Running Your Own Instance

```bash
# Just three steps to get started:
git clone https://github.com/recally-io/recally
cd recally

# Set up your environment (we've added comments to make it clear)
cp env.example .env
# Edit .env with your settings

docker compose up -d
```

Then just open http://localhost:1323 and you're good to go!

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

## License

Free for personal use under MIT. For commercial stuff, drop us a line at [support@recally.io](mailto:support@recally.io).
