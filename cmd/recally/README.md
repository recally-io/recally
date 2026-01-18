# Recally CLI 📚

A standalone command-line tool for saving web articles as markdown files with full metadata preservation. Perfect for archiving articles, building a personal knowledge base, or offline reading.

## ✨ Features

- **🚀 Fast HTTP Mode**: Quick fetching for static pages (default)
- **🌐 Browser Mode**: Handles JavaScript-heavy sites (SPAs, dynamic content)
- **📝 Clean Markdown**: Automatic conversion with readability processing
- **🏷️ Rich Metadata**: YAML frontmatter with title, author, dates, images
- **💾 Smart Storage**: XDG-compliant directory structure with conflict resolution
- **🔍 Verbose Logging**: Optional debug output for troubleshooting

## 📦 Installation

### From Source

**Prerequisites:**
- [Go 1.24+](https://go.dev/dl/)
- [mise](https://mise.jit.su/) (recommended for development)

**Using mise (recommended):**
```bash
# Clone the repository
git clone https://github.com/recally-io/recally
cd recally

# Install mise if not already installed
curl https://mise.run | sh

# Build the CLI
mise run build:cli

# The binary is now at bin/recally
./bin/recally --version
```

**Using Go directly:**
```bash
git clone https://github.com/recally-io/recally
cd recally

# Build with version info
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
go build \
  -ldflags="-s -w -X main.version=$VERSION" \
  -o bin/recally \
  cmd/recally/*.go

./bin/recally --version
```

**Add to PATH (optional):**
```bash
# Copy to a directory in your PATH
sudo cp bin/recally /usr/local/bin/

# Or add bin directory to PATH in ~/.bashrc or ~/.zshrc
export PATH="$PATH:/path/to/recally/bin"
```

### Binary Downloads

> **Coming Soon**: Pre-built binaries will be available in [GitHub Releases](https://github.com/recally-io/recally/releases) for:
> - Linux (amd64, arm64)
> - macOS (amd64, arm64)
> - Windows (planned)

## 🚀 Quick Start

Save your first article:

```bash
# Basic usage (HTTP mode - fastest)
recally https://example.com/article

# Output:
# Fetching https://example.com/article...
# Processing...
# Saved to /home/user/.local/share/recally/contents/2026-01-18/example-article.md
```

For JavaScript-heavy sites:

```bash
# Use browser mode for SPAs and dynamic content
recally --browser https://medium.com/@user/article
```

## 📖 Usage

### HTTP Mode (Default)

Best for static HTML pages and traditional websites. Fast and lightweight.

```bash
recally https://example.com/article
```

**When to use:**
- News articles, blogs, documentation
- Static HTML pages
- Traditional websites without heavy JavaScript

**Advantages:**
- ⚡ Very fast (2-5 seconds typical)
- 💡 Low resource usage
- 🔌 No external dependencies

### Browser Mode

For JavaScript-heavy sites that require full browser rendering.

```bash
recally --browser https://example.com/article
```

**When to use:**
- Single Page Applications (SPAs)
- Sites with lazy-loaded content
- Pages requiring JavaScript execution
- Sites blocking simple HTTP requests

**Advantages:**
- 🌐 Full browser rendering
- 📜 Handles dynamic content
- 🖼️ Captures lazy-loaded images

**Requirements:**
- Chrome/Chromium browser service (see [Browser Service Setup](#browser-service-setup))

### Command-Line Options

```bash
recally [options] <url>

Options:
  --browser              Use browser fetcher for JavaScript-heavy sites
  --browser-url string   Browser control URL (default: "http://localhost:9222")
  --verbose              Enable debug logging
  --output-dir string    Custom output directory (default: XDG data directory)
  --version              Show version information
  -h, --help            Show help message

Positional Arguments:
  url                    Web page URL to fetch (required)
```

### Options Reference

#### `--browser`

Enable browser mode for JavaScript-heavy sites.

```bash
recally --browser https://medium.com/@user/article
```

#### `--browser-url <url>`

Specify custom Chrome DevTools Protocol control URL.

```bash
recally --browser --browser-url http://remote-browser:9222 https://example.com
```

**Default:** `http://localhost:9222`  
**Environment Variable:** `BROWSER_CONTROL_URL` (flag takes precedence)

#### `--verbose`

Enable debug logging for troubleshooting.

```bash
recally --verbose https://example.com/article
```

**Output includes:**
- Configuration details
- Timing information
- Content processing stats
- File operations

**Example output:**
```
level=INFO msg="recally configuration" version=v1.0.0 url=https://example.com browser_mode=false
level=INFO msg="output directory created" path=/home/user/.local/share/recally/contents/2026-01-18
level=INFO msg="fetcher created" mode=HTTP
Fetching https://example.com/article...
Processing...
level=INFO msg="content fetched and processed" title="Example Article" content_length=5234 elapsed=2.3s
Saved to /home/user/.local/share/recally/contents/2026-01-18/example-article.md
level=INFO msg="operation completed successfully" total_elapsed=2.5s
```

#### `--output-dir <path>`

Save articles to a custom directory instead of XDG default.

```bash
recally --output-dir ~/my-articles https://example.com/article
```

**Default behavior:**
- Linux: `~/.local/share/recally/contents/YYYY-MM-DD/`
- macOS: `~/Library/Application Support/recally/contents/YYYY-MM-DD/`
- Windows: `%LOCALAPPDATA%\recally\contents\YYYY-MM-DD\`

**Custom directory structure:**
```
~/my-articles/
└── YYYY-MM-DD/
    ├── article-1.md
    └── article-2.md
```

## ⚙️ Configuration

### Environment Variables

#### `BROWSER_CONTROL_URL`

Default browser control URL for browser mode.

```bash
export BROWSER_CONTROL_URL=http://localhost:9222
recally --browser https://example.com/article
```

**Precedence:** CLI flag (`--browser-url`) > Environment variable > Default value

**Example with remote browser:**
```bash
export BROWSER_CONTROL_URL=http://192.168.1.100:9222
recally --browser https://example.com/article
```

### Output File Format

Articles are saved as markdown files with YAML frontmatter:

```markdown
---
url: https://example.com/article
title: Example Article Title
author: John Doe
description: A comprehensive guide to example topics
site_name: Example Site
published_time: 2026-01-15T10:00:00Z
modified_time: 2026-01-17T12:30:00Z
cover: https://example.com/images/cover.jpg
favicon: https://example.com/favicon.ico
saved_at: 2026-01-18T08:30:00Z
---

# Example Article Title

Article content in clean markdown format...
```

**Frontmatter fields:**
- `url`: Original article URL (required)
- `title`: Article title (required)
- `author`: Article author (if available)
- `description`: Article description/excerpt (if available)
- `site_name`: Website name (if available)
- `published_time`: Original publication time in RFC3339 format (if available)
- `modified_time`: Last modified time in RFC3339 format (if available)
- `cover`: Cover image URL (if available)
- `favicon`: Site favicon URL (if available)
- `saved_at`: Timestamp when article was saved (UTC, RFC3339 format)

### File Naming

Articles are automatically saved with sanitized filenames based on the title:

**Sanitization rules:**
1. Trim whitespace and convert to lowercase
2. Replace multiple spaces with single hyphen
3. Remove non-alphanumeric characters (except hyphens)
4. Preserve Unicode letters and numbers (supports international characters)
5. Truncate to 200 characters (adds MD5 hash suffix if truncated)
6. Fallback to `untitled-{timestamp}` if title is empty

**Examples:**
| Original Title | Filename |
|----------------|----------|
| "Getting Started with Go" | `getting-started-with-go.md` |
| "React 18: What's New?" | `react-18-whats-new.md` |
| "文章标题" | `文章标题.md` |
| (no title) | `untitled-1737187200.md` |

**Conflict resolution:**

If a file with the same name exists, a counter is appended:
- First save: `article.md`
- Conflict 1: `article-1.md`
- Conflict 2: `article-2.md`

### Browser Service Setup

Browser mode requires a Chrome/Chromium instance with DevTools Protocol enabled.

#### Docker Setup (Recommended)

**Using Docker Compose:**

Create `docker-compose.yml`:
```yaml
version: '3.8'

services:
  chrome:
    image: chromedp/headless-shell:latest
    ports:
      - "9222:9222"
    restart: unless-stopped
    environment:
      # Optional: Configure Chrome flags
      - CHROME_ARGS=--no-sandbox --disable-gpu
    # Optional: Increase shared memory for complex pages
    shm_size: 2gb
```

**Start the service:**
```bash
docker compose up -d

# Verify it's running
curl -s http://localhost:9222/json/version | jq .
```

**Using Docker directly:**
```bash
docker run -d \
  --name recally-chrome \
  -p 9222:9222 \
  --shm-size=2gb \
  chromedp/headless-shell:latest
```

#### Local Browser Setup

**Chrome/Chromium:**
```bash
# Linux
google-chrome --headless --remote-debugging-port=9222

# macOS
/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome \
  --headless --remote-debugging-port=9222

# With custom user data directory (recommended)
google-chrome \
  --headless \
  --remote-debugging-port=9222 \
  --user-data-dir=/tmp/chrome-recally
```

**Verify connection:**
```bash
# Check if browser control is accessible
curl http://localhost:9222/json/version

# Test with recally
recally --browser --verbose https://example.com/article
```

#### Remote Browser Setup

For running browser on a different machine:

**On browser host:**
```bash
docker run -d \
  --name recally-chrome \
  -p 9222:9222 \
  chromedp/headless-shell:latest
```

**On recally host:**
```bash
export BROWSER_CONTROL_URL=http://browser-host:9222
recally --browser https://example.com/article

# Or use flag
recally --browser --browser-url http://browser-host:9222 https://example.com
```

## 💡 Examples

### Basic Article Saving

```bash
# Save a news article
recally https://news.ycombinator.com/item?id=12345

# Save a blog post
recally https://blog.example.com/2026/01/great-post

# Save documentation
recally https://docs.example.com/guide/getting-started
```

### Browser Mode for Dynamic Content

```bash
# Medium article (requires JavaScript)
recally --browser https://medium.com/@user/article

# Modern web app content
recally --browser https://app.example.com/article

# SPA with lazy loading
recally --browser https://spa.example.com/content/123
```

### Custom Output Location

```bash
# Save to specific directory
recally --output-dir ~/Documents/articles https://example.com/article

# Organize by topic
recally --output-dir ~/articles/tech https://tech-blog.com/post
recally --output-dir ~/articles/science https://science-news.com/article
```

### Remote Browser

```bash
# Using environment variable
export BROWSER_CONTROL_URL=http://192.168.1.100:9222
recally --browser https://example.com/article

# Using flag
recally --browser --browser-url http://192.168.1.100:9222 https://example.com
```

### Debugging and Troubleshooting

```bash
# Enable verbose logging
recally --verbose https://example.com/article

# Debug browser mode
recally --browser --verbose https://difficult-site.com/article

# Test with known-good URL
recally --verbose https://example.com
```

### Batch Processing (Shell Script)

```bash
#!/bin/bash
# save-articles.sh - Batch save multiple articles

urls=(
  "https://blog.example.com/post-1"
  "https://blog.example.com/post-2"
  "https://blog.example.com/post-3"
)

for url in "${urls[@]}"; do
  echo "Saving: $url"
  recally "$url"
  sleep 2  # Be nice to the server
done

echo "All articles saved!"
```

### Integration with Other Tools

**Save from clipboard:**
```bash
# Linux (xclip)
recally "$(xclip -o -selection clipboard)"

# macOS
recally "$(pbpaste)"
```

**Save with notification:**
```bash
#!/bin/bash
url="$1"
if recally "$url"; then
  notify-send "Recally" "Article saved successfully"
else
  notify-send "Recally" "Failed to save article" --urgency=critical
fi
```

**Integration with browser bookmarklet:**

Save this as a browser bookmarklet:
```javascript
javascript:(function(){window.open('http://your-server.com/save?url='+encodeURIComponent(location.href))})()
```

Then create a simple web service that calls recally:
```bash
# save-server.sh
#!/bin/bash
URL=$(echo "$QUERY_STRING" | sed -n 's/^.*url=\([^&]*\).*$/\1/p' | urldecode)
recally "$URL"
```

## 🔧 Troubleshooting

### Common Errors

#### "Error: URL is required"

**Cause:** No URL was provided.

**Solution:**
```bash
# ❌ Wrong
recally

# ✅ Correct
recally https://example.com/article
```

#### "Error: Invalid URL: URL must include http:// or https:// scheme"

**Cause:** URL is missing the protocol.

**Solution:**
```bash
# ❌ Wrong
recally example.com/article

# ✅ Correct
recally https://example.com/article
```

#### "Error: Failed to fetch and process content: failed to fetch: context deadline exceeded"

**Cause:** Network timeout (operation took longer than 5 minutes).

**Solutions:**
- Check your internet connection
- Try again (site might be temporarily slow)
- Use browser mode if HTTP mode is timing out: `recally --browser URL`
- Check if the site is accessible: `curl -I URL`

#### "Error: Failed to create fetcher: create browser fetcher: failed to connect"

**Cause:** Browser service is not running or not accessible.

**Solutions:**

1. **Start the browser service:**
   ```bash
   # Using Docker
   docker run -d -p 9222:9222 chromedp/headless-shell:latest
   
   # Using local Chrome
   google-chrome --headless --remote-debugging-port=9222
   ```

2. **Verify browser is running:**
   ```bash
   curl http://localhost:9222/json/version
   # Should return JSON with browser version
   ```

3. **Check if port is in use:**
   ```bash
   # Linux/macOS
   lsof -i :9222
   
   # If different port, use --browser-url flag
   recally --browser --browser-url http://localhost:9223 URL
   ```

4. **Check Docker logs:**
   ```bash
   docker logs recally-chrome
   ```

#### "Error: Filesystem error: permission denied"

**Cause:** No write permission to output directory.

**Solutions:**

1. **Check directory permissions:**
   ```bash
   ls -ld ~/.local/share/recally
   ```

2. **Fix permissions:**
   ```bash
   chmod 755 ~/.local/share/recally
   ```

3. **Use custom output directory:**
   ```bash
   recally --output-dir ~/my-articles URL
   ```

#### "Error: Filesystem error: insufficient disk space"

**Cause:** Less than 100MB of free disk space.

**Solutions:**

1. **Check disk usage:**
   ```bash
   df -h ~/.local/share/recally
   ```

2. **Free up space:**
   ```bash
   # Find large files
   du -sh ~/.local/share/recally/* | sort -h
   
   # Remove old articles if needed
   rm -rf ~/.local/share/recally/contents/2025-*
   ```

3. **Use different disk:**
   ```bash
   recally --output-dir /mnt/other-disk/articles URL
   ```

### Network Issues

#### SSL/TLS Certificate Errors

**Error:** `x509: certificate signed by unknown authority`

**Cause:** Site uses self-signed certificate or corporate proxy.

**Solution:**
- Use browser mode (handles certificates better)
- If on corporate network, check with IT about proxy configuration

#### Connection Refused

**Error:** `connection refused`

**Cause:** Site is blocking requests or is down.

**Solutions:**
- Verify site is accessible in a browser
- Use browser mode: `recally --browser URL`
- Check if site requires specific headers or cookies

#### Rate Limiting

**Error:** `429 Too Many Requests` or similar

**Cause:** Site is rate-limiting requests.

**Solutions:**
- Wait and retry later
- Use browser mode (appears more like real user)
- Add delays between requests in scripts

### Browser Service Issues

#### Port Already in Use

**Error:** `bind: address already in use`

**Cause:** Port 9222 is already occupied.

**Solutions:**

1. **Find process using port:**
   ```bash
   lsof -i :9222
   ```

2. **Stop existing process:**
   ```bash
   # If it's Docker
   docker stop recally-chrome
   
   # If it's local Chrome
   pkill -f "chrome.*remote-debugging-port=9222"
   ```

3. **Use different port:**
   ```bash
   # Start browser on different port
   docker run -d -p 9223:9222 chromedp/headless-shell:latest
   
   # Use custom port
   recally --browser --browser-url http://localhost:9223 URL
   ```

#### Browser Out of Memory

**Symptoms:** Browser crashes or hangs on complex pages.

**Solutions:**

1. **Increase shared memory (Docker):**
   ```bash
   docker run -d \
     -p 9222:9222 \
     --shm-size=2gb \
     chromedp/headless-shell:latest
   ```

2. **Restart browser service:**
   ```bash
   docker restart recally-chrome
   ```

3. **Use HTTP mode if possible:**
   ```bash
   recally URL  # Try without --browser first
   ```

### Performance Issues

#### Slow Fetching

**Issue:** Articles take long time to fetch.

**Solutions:**
- Use HTTP mode (faster than browser mode)
- Check internet connection speed
- Try at different time (site might be slow)
- Enable verbose mode to see where time is spent: `recally --verbose URL`

#### Large Binary Size

**Issue:** recally binary is large.

**Note:** Binary is optimized with `-ldflags="-s -w"` (strip debug info). Typical size is 15-25MB depending on platform.

**Optional:** Further compress with upx (not recommended - may cause issues):
```bash
upx --best bin/recally
```

### Getting Help

**Enable verbose logging:**
```bash
recally --verbose URL 2>&1 | tee recally.log
```

**Check version:**
```bash
recally --version
```

**File an issue:**

If you encounter a bug, please open an issue at:
https://github.com/recally-io/recally/issues

Include:
- Recally version (`recally --version`)
- Operating system and version
- Full command you ran
- Error output (with `--verbose` if possible)
- Example URL (if not private)

## 📊 Exit Codes

The CLI uses standard exit codes for scripting and automation:

| Code | Meaning | Common Causes |
|------|---------|---------------|
| **0** | Success | Article saved successfully |
| **1** | Fetch/Process Error | Network timeout, parsing failure, invalid HTML, browser connection failed |
| **2** | Usage Error | Invalid flags, missing URL, malformed URL, wrong number of arguments |
| **3** | Filesystem Error | Permission denied, disk full, read-only filesystem, directory creation failed |

**Example usage in scripts:**
```bash
#!/bin/bash
recally "$url"
exit_code=$?

case $exit_code in
  0)
    echo "✅ Success"
    ;;
  1)
    echo "❌ Network or parsing error"
    exit 1
    ;;
  2)
    echo "❌ Invalid command"
    exit 2
    ;;
  3)
    echo "❌ Filesystem error (check permissions and disk space)"
    exit 3
    ;;
esac
```

## 🔒 Security Considerations

### URL Validation

- Only `http://` and `https://` schemes are allowed
- `file://`, `javascript:`, and `data:` URLs are rejected
- URLs are validated before fetching

### Filename Sanitization

- All filenames are sanitized to prevent directory traversal
- Special characters are removed or replaced with safe equivalents
- Maximum filename length is enforced (200 chars + hash)

### Disk Space Checks

- Requires at least 100MB free space before writing
- Prevents disk full errors and system instability

### Symlink Protection

- Directory creation uses `os.MkdirAll` which is safe against symlink attacks
- No following of symbolic links during path resolution

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. **Fork and clone:**
   ```bash
   git clone https://github.com/YOUR-USERNAME/recally
   cd recally
   ```

2. **Install dependencies:**
   ```bash
   mise install
   ```

3. **Make your changes:**
   - CLI code is in `cmd/recally/`
   - Add tests for new features
   - Follow existing code style

4. **Test your changes:**
   ```bash
   # Run tests
   go test ./cmd/recally/...
   
   # Build and test manually
   mise run build:cli
   ./bin/recally --version
   ./bin/recally --verbose https://example.com
   ```

5. **Submit a pull request:**
   - Write clear commit messages
   - Update documentation if needed
   - Reference any related issues

See [CONTRIBUTING.md](../../CONTRIBUTING.md) for more details.

## 📄 License

Recally is licensed under:
- **Non-commercial use**: [AGPLv3](../../LICENSE)
- **Commercial use**: Contact [sales@recally.io](mailto:sales@recally.io)

## 🔗 Related Projects

**Main Recally Application:**
- [Recally Web App](../../README.md) - Full-featured web application with AI features
- [Browser Extensions](../../extensions/) - Chrome and Firefox extensions
- [Telegram Bot](../../internal/port/bots/) - Save articles via Telegram

**Alternative Tools:**
- [Shiori](https://github.com/go-shiori/shiori) - Bookmark manager
- [Omnivore](https://github.com/omnivore-app/omnivore) - Read-it-later app
- [Hoarder](https://github.com/hoarder-app/hoarder) - Self-hosted bookmark manager

## 🙏 Acknowledgments

Built with:
- [go-shiori/go-readability](https://github.com/go-shiori/go-readability) - Content extraction
- [JohannesKaufmann/html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) - HTML to Markdown conversion
- [go-rod/rod](https://github.com/go-rod/rod) - Browser automation
- [yaml.v3](https://github.com/go-yaml/yaml) - YAML processing

---

**Made with ❤️ by the Recally Team**

For questions or support, reach out:
- 📧 Email: [support@recally.io](mailto:support@recally.io)
- 💬 GitHub Issues: [github.com/recally-io/recally/issues](https://github.com/recally-io/recally/issues)
- 🌐 Website: [recally.io](https://recally.io)
