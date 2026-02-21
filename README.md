# prowl4ai

`prowl4ai` is a lightweight Go web crawler inspired by [crawl4ai](https://github.com/unclecode/crawl4ai).

It fetches fully rendered HTML with Playwright, runs readability-style cleaning, converts content to Markdown, and returns structured crawl results.

## Features

- Playwright-backed page rendering (JavaScript-enabled sites)
- Configurable crawl timeout and headless mode
- Readability-based clean HTML extraction
- Markdown conversion from cleaned HTML
- JSON or Markdown CLI output
- Crawl metadata (title, byline, excerpt, language)
- Response metadata (status code, headers, redirected URL)

## Project Status

This is an early-stage project with a single CLI command: `crawl`.

Current CLI flags intentionally expose a minimal surface area:

- `--timeout` (page timeout in milliseconds)
- `--headless` (run browser headless or headed)
- `--output` (`json` or `markdown`)

Additional crawler options exist in internal config types and can be exposed as the CLI evolves.

## Requirements

- Go `1.25+`
- Playwright runtime + browser binaries

## Getting Started

1. Clone and enter the repo:

```bash
git clone <your-fork-or-repo-url>
cd prowl4ai
```

2. Install Go dependencies:

```bash
go mod download
```

3. Install Playwright browser binaries (Chromium):

```bash
go run github.com/playwright-community/playwright-go/cmd/playwright@latest install chromium
```

## Usage

Run a crawl and return JSON:

```bash
go run ./cmd/prowl4ai crawl https://example.com
```

Run with custom timeout and headed browser:

```bash
go run ./cmd/prowl4ai crawl --timeout 90000 --headless=false https://example.com
```

Return Markdown only:

```bash
go run ./cmd/prowl4ai crawl --output markdown https://example.com
```

### CLI Help

```text
Usage:
  prowl4ai crawl [--timeout ms] [--headless bool] [--output json|markdown] <url>
```

## Use as a Go Library

Install the module:

```bash
go get github.com/techbysteve/prowl4ai
```

Then use it directly in code:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/techbysteve/prowl4ai"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	crawler := prowl4ai.NewCrawler()
	defer func() {
		if err := crawler.Close(context.Background()); err != nil {
			log.Printf("close error: %v", err)
		}
	}()

	result, err := crawler.Crawl(ctx, "https://example.com")
	if err != nil {
		log.Fatalf("crawl failed: %v", err)
	}

	fmt.Println(result.Markdown)
}
```

Custom browser/run config:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/techbysteve/prowl4ai"
)

func main() {
	browserCfg := prowl4ai.DefaultBrowserConfig()
	browserCfg.Headless = true

	runCfg := prowl4ai.DefaultRunConfig()
	runCfg.PageTimeoutMs = 120000
	runCfg.WaitUntil = "networkidle"
	runCfg.WaitFor = "main article"
	runCfg.WaitForTimeoutMs = 10000

	crawler := prowl4ai.NewCrawlerWithConfig(browserCfg, runCfg)
	defer func() {
		if err := crawler.Close(context.Background()); err != nil {
			log.Printf("close error: %v", err)
		}
	}()

	result, err := crawler.Crawl(context.Background(), "https://example.com")
	if err != nil {
		log.Fatalf("crawl failed: %v", err)
	}

	fmt.Printf("status=%d redirected=%s\n", result.StatusCode, result.RedirectedURL)
}
```

## Output Shape (JSON)

A successful crawl returns fields like:

- `url`
- `html`
- `cleaned_html`
- `markdown`
- `metadata`
- `status_code`
- `response_headers`
- `redirected_url`
- `success`

On failure, the tool still returns structured JSON with:

- `success: false`
- `error_message`
- best-effort context fields when available

## Repository Layout

- `cmd/prowl4ai/main.go`: CLI entrypoint
- `internal/browser/`: browser adapter abstraction + Playwright implementation
- `internal/prowler/`: crawler service orchestration
- `internal/extract/`: clean HTML + Markdown pipeline
- `internal/config/`: browser and run defaults
- `internal/model/`: crawl result models

## Notes and Limitations

- Browser backend is Playwright-based and defaults to Chromium.
- CLI currently exposes only a subset of internal crawl configuration.
- Link extraction is not yet exposed in CLI output.

## Inspiration

This project is inspired by the ideas and workflow of crawl4ai, adapted into a focused Go implementation.
