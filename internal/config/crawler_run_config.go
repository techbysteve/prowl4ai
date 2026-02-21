package config

const (
	DefaultPageTimeoutMs = 60000
	DefaultWaitUntil     = "domcontentloaded"
)

// CrawlerRunConfig controls a single crawl execution.
// Keep this small and stable for Phase 1; extend in later phases as needed.
type CrawlerRunConfig struct {
	PageTimeoutMs    int    `json:"page_timeout_ms"`
	WaitUntil        string `json:"wait_until"`
	WaitFor          string `json:"wait_for,omitempty"`
	WaitForTimeoutMs int    `json:"wait_for_timeout_ms,omitempty"`
	EnableCleanHTML  bool   `json:"enable_clean_html"`
	EnableMarkdown   bool   `json:"enable_markdown"`
	EnableLinks      bool   `json:"enable_links"`
	OnlyText         bool   `json:"only_text"`
	CSSSelector      string `json:"css_selector,omitempty"`
	Verbose          bool   `json:"verbose"`
}

func DefaultCrawlerRunConfig() CrawlerRunConfig {
	return CrawlerRunConfig{
		PageTimeoutMs:    DefaultPageTimeoutMs,
		WaitUntil:        DefaultWaitUntil,
		WaitFor:          "",
		WaitForTimeoutMs: 0,
		EnableCleanHTML:  true,
		EnableMarkdown:   true,
		EnableLinks:      false,
		OnlyText:         false,
		CSSSelector:      "",
		Verbose:          true,
	}
}
