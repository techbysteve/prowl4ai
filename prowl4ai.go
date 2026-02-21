package prowl4ai

import (
	"context"

	"github.com/techbysteve/prowl4ai/internal/browser"
	"github.com/techbysteve/prowl4ai/internal/config"
	"github.com/techbysteve/prowl4ai/internal/model"
	"github.com/techbysteve/prowl4ai/internal/prowler"
)

// CrawlResult is the structured output returned by a crawl run.
type CrawlResult = model.CrawlResult

// BrowserConfig controls browser startup behavior for library users.
type BrowserConfig struct {
	BrowserType    string
	Headless       bool
	Channel        string
	Proxy          string
	ViewportWidth  int
	ViewportHeight int
	UserAgent      string
	ExtraArgs      []string
	DebuggingPort  int
}

// RunConfig controls a single crawl execution.
type RunConfig struct {
	PageTimeoutMs    int
	WaitUntil        string
	WaitFor          string
	WaitForTimeoutMs int
	EnableCleanHTML  bool
	EnableMarkdown   bool
	OnlyText         bool
	CSSSelector      string
	Verbose          bool
}

// DefaultBrowserConfig returns sensible browser defaults.
func DefaultBrowserConfig() BrowserConfig {
	cfg := config.DefaultBrowserConfig()
	return BrowserConfig{
		BrowserType:    cfg.BrowserType,
		Headless:       cfg.Headless,
		Channel:        cfg.Channel,
		Proxy:          cfg.Proxy,
		ViewportWidth:  cfg.ViewportWidth,
		ViewportHeight: cfg.ViewportHeight,
		UserAgent:      cfg.UserAgent,
		ExtraArgs:      append([]string{}, cfg.ExtraArgs...),
		DebuggingPort:  cfg.DebuggingPort,
	}
}

// DefaultRunConfig returns sensible crawl defaults.
func DefaultRunConfig() RunConfig {
	cfg := config.DefaultCrawlerRunConfig()
	return RunConfig{
		PageTimeoutMs:    cfg.PageTimeoutMs,
		WaitUntil:        cfg.WaitUntil,
		WaitFor:          cfg.WaitFor,
		WaitForTimeoutMs: cfg.WaitForTimeoutMs,
		EnableCleanHTML:  cfg.EnableCleanHTML,
		EnableMarkdown:   cfg.EnableMarkdown,
		OnlyText:         cfg.OnlyText,
		CSSSelector:      cfg.CSSSelector,
		Verbose:          cfg.Verbose,
	}
}

// Crawler is the main Go library entrypoint.
type Crawler struct {
	service          *prowler.Service
	defaultRunConfig config.CrawlerRunConfig
}

// NewCrawler builds a crawler with default browser and run configuration.
func NewCrawler() *Crawler {
	return NewCrawlerWithConfig(DefaultBrowserConfig(), DefaultRunConfig())
}

// NewCrawlerWithConfig builds a crawler with explicit browser and run defaults.
func NewCrawlerWithConfig(browserCfg BrowserConfig, runCfg RunConfig) *Crawler {
	internalBrowserCfg := toInternalBrowserConfig(browserCfg)
	adapter := browser.NewPlaywrightAdapter(internalBrowserCfg)
	service := prowler.NewService(adapter)
	return &Crawler{
		service:          service,
		defaultRunConfig: toInternalRunConfig(runCfg),
	}
}

// SetDefaultRunConfig updates the default run config used by Crawl.
func (c *Crawler) SetDefaultRunConfig(cfg RunConfig) {
	c.defaultRunConfig = toInternalRunConfig(cfg)
}

// Crawl executes one crawl using the crawler's default run config.
func (c *Crawler) Crawl(ctx context.Context, url string) (CrawlResult, error) {
	return c.service.Run(ctx, url, c.defaultRunConfig)
}

// CrawlWithConfig executes one crawl using a per-call run config.
func (c *Crawler) CrawlWithConfig(ctx context.Context, url string, cfg RunConfig) (CrawlResult, error) {
	return c.service.Run(ctx, url, toInternalRunConfig(cfg))
}

// Close releases browser resources.
func (c *Crawler) Close(ctx context.Context) error {
	return c.service.Close(ctx)
}

func toInternalBrowserConfig(cfg BrowserConfig) config.BrowserConfig {
	base := config.DefaultBrowserConfig()
	base.BrowserType = cfg.BrowserType
	base.Headless = cfg.Headless
	base.Channel = cfg.Channel
	base.Proxy = cfg.Proxy
	base.ViewportWidth = cfg.ViewportWidth
	base.ViewportHeight = cfg.ViewportHeight
	base.UserAgent = cfg.UserAgent
	base.ExtraArgs = append([]string{}, cfg.ExtraArgs...)
	base.DebuggingPort = cfg.DebuggingPort
	return base
}

func toInternalRunConfig(cfg RunConfig) config.CrawlerRunConfig {
	base := config.DefaultCrawlerRunConfig()
	base.PageTimeoutMs = cfg.PageTimeoutMs
	base.WaitUntil = cfg.WaitUntil
	base.WaitFor = cfg.WaitFor
	base.WaitForTimeoutMs = cfg.WaitForTimeoutMs
	base.EnableCleanHTML = cfg.EnableCleanHTML
	base.EnableMarkdown = cfg.EnableMarkdown
	base.OnlyText = cfg.OnlyText
	base.CSSSelector = cfg.CSSSelector
	base.Verbose = cfg.Verbose
	return base
}
