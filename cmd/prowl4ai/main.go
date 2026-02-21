package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/techbysteve/prowl4ai/internal/browser"
	"github.com/techbysteve/prowl4ai/internal/config"
	"github.com/techbysteve/prowl4ai/internal/prowler"
)

func main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) < 2 {
		printUsage()
		return 2
	}

	switch os.Args[1] {
	case "crawl":
		return runCrawl(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		return 2
	}
}

func runCrawl(args []string) int {
	fs := flag.NewFlagSet("crawl", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	timeoutMs := fs.Int("timeout", config.DefaultPageTimeoutMs, "Page timeout in milliseconds")
	headless := fs.Bool("headless", true, "Run browser in headless mode")
	output := fs.String("output", "json", "Output format: json|markdown")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: prowl4ai crawl [--timeout ms] [--headless bool] [--output json|markdown] <url>")
		return 2
	}
	url := fs.Arg(0)
	if *output != "json" && *output != "markdown" {
		fmt.Fprintln(os.Stderr, "invalid --output value, expected: json|markdown")
		return 2
	}

	browserCfg := config.DefaultBrowserConfig()
	browserCfg.Headless = *headless

	runCfg := config.DefaultCrawlerRunConfig()
	runCfg.PageTimeoutMs = *timeoutMs

	adapter := browser.NewPlaywrightAdapter(browserCfg)
	service := prowler.NewService(adapter)

	ctx := context.Background()
	defer func() {
		if err := service.Close(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "close error: %v\n", err)
		}
	}()

	result, err := service.Run(ctx, url, runCfg)
	if err != nil {
		// Still print structured result to aid debugging.
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(result)
		return 1
	}

	if *output == "markdown" {
		if result.Markdown != "" {
			fmt.Print(result.Markdown)
		} else if result.CleanedHTML != "" {
			fmt.Print(result.CleanedHTML)
		} else {
			fmt.Print(result.HTML)
		}
		return 0
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode result: %v\n", err)
		return 1
	}
	return 0
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  prowl4ai crawl [--timeout ms] [--headless bool] [--output json|markdown] <url>")
}
