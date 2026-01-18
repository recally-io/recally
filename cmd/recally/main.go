package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

// version will be injected at build time via ldflags
var version = "dev"

// Exit codes
const (
	ExitSuccess     = 0 // Success
	ExitFetchError  = 1 // Fetch/process error (network, parsing failures)
	ExitUsageError  = 2 // Usage error (invalid flags, missing URL)
	ExitFSError     = 3 // Filesystem error (permissions, disk full)
)

// CLI flags
var (
	flagBrowser    bool
	flagBrowserURL string
	flagVerbose    bool
	flagOutputDir  string
	flagVersion    bool
)

func init() {
	// Define flags
	flag.BoolVar(&flagBrowser, "browser", false, "Use browser fetcher instead of HTTP")
	flag.StringVar(&flagBrowserURL, "browser-url", "", "Browser control URL (empty = launch new browser)")
	flag.BoolVar(&flagVerbose, "verbose", false, "Enable debug logging")
	flag.StringVar(&flagOutputDir, "output-dir", "", "Custom output directory (empty = use XDG default)")
	flag.BoolVar(&flagVersion, "version", false, "Show version information")

	// Customize usage message
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `recally - Save web articles as markdown

Usage:
  recally [options] <url>

Options:
  --browser              Use browser fetcher for JavaScript-heavy sites (default: false)
  --browser-url string   Browser control URL (empty = launch new browser)
  --verbose              Enable debug logging (default: false)
  --output-dir string    Custom output directory (empty = use XDG default)
  --version              Show version information
  -h, --help             Show this help message

Examples:
  # Basic usage with HTTP fetcher
  recally https://example.com/article

  # Use browser fetcher (launches new Chrome instance)
  recally --browser https://example.com/article

  # Use existing browser service
  recally --browser --browser-url http://localhost:9222 https://example.com/article

  # Enable verbose logging
  recally --verbose https://example.com/article

  # Custom output directory
  recally --output-dir ~/my-articles https://example.com/article

Exit Codes:
  0 - Success
  1 - Fetch/process error (network, parsing failures)
  2 - Usage error (invalid flags, missing URL)
  3 - Filesystem error (permissions, disk full)

Environment Variables:
  BROWSER_CONTROL_URL    Browser control URL (overridden by --browser-url flag)
`)
}

func main() {
	os.Exit(run())
}

func run() int {
	// Parse flags
	flag.Parse()

	// Show version and exit
	if flagVersion {
		fmt.Printf("recally version %s\n", version)
		return ExitSuccess
	}

	// Get URL from positional argument
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: URL is required")
		fmt.Fprintln(os.Stderr, "Run 'recally --help' for usage information")
		return ExitUsageError
	}

	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Error: Only one URL is allowed")
		fmt.Fprintln(os.Stderr, "Run 'recally --help' for usage information")
		return ExitUsageError
	}

	articleURL := args[0]

	// Validate URL
	if err := validateURL(articleURL); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid URL: %v\n", err)
		return ExitUsageError
	}

	// Override browser URL from environment if not set via flag
	if flagBrowser && flagBrowserURL == "" {
		if envURL := os.Getenv("BROWSER_CONTROL_URL"); envURL != "" {
			flagBrowserURL = envURL
		}
	}

	// Execute main workflow
	return execute(articleURL)
}

// execute runs the main workflow: fetch, process, and save content
func execute(articleURL string) int {
	// Track start time for elapsed time calculation
	startTime := time.Now()

	// 1. Create logger
	logger := NewLogger(flagVerbose)

	// Verbose mode: Show configuration
	if flagVerbose {
		logger.Info("recally configuration",
			"version", version,
			"url", articleURL,
			"browser_mode", flagBrowser,
			"browser_url", getBrowserURLDescription(),
			"output_dir", getOutputDirDescription(),
		)
	}

	// 2. Create context with 5-minute hard timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 3. Get output directory
	outputDir, err := GetOutputDir(flagOutputDir, time.Now())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create output directory: %v\n", err)
		return ExitFSError
	}

	if flagVerbose {
		logger.Info("output directory created", "path", outputDir)
	}

	// 4. Create fetcher
	fetcher, err := NewFetcher(flagBrowser, flagBrowserURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create fetcher: %v\n", err)
		return ExitFetchError
	}

	if flagVerbose {
		mode := "HTTP"
		if flagBrowser {
			mode = "Browser"
		}
		logger.Info("fetcher created", "mode", mode)
	}

	// 5. Fetch and process content
	fmt.Fprintf(os.Stderr, "Fetching %s...\n", articleURL)

	content, err := FetchAndProcess(ctx, articleURL, logger, fetcher)
	if err != nil {
		// Check if error is due to context cancellation
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Fprintf(os.Stderr, "Error: Operation timed out after 5 minutes\n")
			return ExitFetchError
		}
		if ctx.Err() == context.Canceled {
			fmt.Fprintf(os.Stderr, "Error: Operation was canceled\n")
			return ExitFetchError
		}

		// Regular fetch/process error
		fmt.Fprintf(os.Stderr, "Error: Failed to fetch and process content: %v\n", err)
		return ExitFetchError
	}

	fmt.Fprintf(os.Stderr, "Processing...\n")

	if flagVerbose {
		elapsed := time.Since(startTime)
		logger.Info("content fetched and processed",
			"title", content.Title,
			"author", content.Author,
			"content_length", len(content.Markwdown),
			"elapsed", elapsed.Round(time.Millisecond).String(),
		)
	}

	// 6. Write markdown to disk
	outputPath, err := WriteMarkdown(content, outputDir)
	if err != nil {
		// Determine if this is a filesystem error
		if isFilesystemError(err) {
			fmt.Fprintf(os.Stderr, "Error: Filesystem error: %v\n", err)
			return ExitFSError
		}

		// Otherwise, treat as general processing error
		fmt.Fprintf(os.Stderr, "Error: Failed to write markdown: %v\n", err)
		return ExitFetchError
	}

	// 7. Success! Print confirmation
	elapsed := time.Since(startTime)
	fmt.Fprintf(os.Stderr, "Saved to %s\n", outputPath)

	if flagVerbose {
		logger.Info("operation completed successfully",
			"output_path", outputPath,
			"total_elapsed", elapsed.Round(time.Millisecond).String(),
		)
	}

	return ExitSuccess
}

// validateURL validates that the URL is properly formatted and uses http/https scheme
func validateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Parse URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	// Check scheme
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		if scheme == "" {
			return fmt.Errorf("URL must include http:// or https:// scheme")
		}
		return fmt.Errorf("unsupported URL scheme '%s' (only http and https are allowed)", scheme)
	}

	// Check host
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must include a host")
	}

	return nil
}

// getOutputDirDescription returns a human-readable description of the output directory
func getOutputDirDescription() string {
	if flagOutputDir != "" {
		return flagOutputDir
	}
	return "(using XDG default)"
}

// isFilesystemError checks if an error is filesystem-related
// Returns true for: permission denied, disk full, no space, read-only filesystem
func isFilesystemError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	// Common filesystem error patterns
	fsErrors := []string{
		"permission denied",
		"disk full",
		"no space left",
		"read-only file system",
		"insufficient disk space",
		"failed to create directory",
		"failed to write file",
		"disk space check failed",
	}

	for _, pattern := range fsErrors {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	return false
}

// getBrowserURLDescription returns a human-readable description of the browser URL
func getBrowserURLDescription() string {
	if flagBrowserURL != "" {
		return flagBrowserURL
	}
	return "(launching new browser)"
}
