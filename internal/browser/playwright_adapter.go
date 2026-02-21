package browser

import (
	"context"
	"fmt"
	"sync"

	"github.com/playwright-community/playwright-go"
	"github.com/techbysteve/prowl4ai/internal/config"
	"github.com/techbysteve/prowl4ai/internal/stderrors"
)

type PlaywrightAdapter struct {
	cfg     config.BrowserConfig
	mu      sync.Mutex
	ready   bool
	pw      *playwright.Playwright
	browser playwright.Browser
}

func NewPlaywrightAdapter(cfg config.BrowserConfig) *PlaywrightAdapter {
	return &PlaywrightAdapter{
		cfg: cfg,
	}
}

func (a *PlaywrightAdapter) Start(ctx context.Context) error {
	_ = ctx
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.ready {
		return nil
	}
	pw, err := playwright.Run()
	if err != nil {
		return err
	}

	launchOptions := a.buildLaunchOptions()
	var browser playwright.Browser
	switch a.cfg.BrowserType {
	case "", "chromium":
		browser, err = pw.Chromium.Launch(launchOptions)
	case "firefox":
		browser, err = pw.Firefox.Launch(launchOptions)
	case "webkit":
		browser, err = pw.WebKit.Launch(launchOptions)
	default:
		_ = pw.Stop()
		return fmt.Errorf("unsupported browser type: %s", a.cfg.BrowserType)
	}
	if err != nil {
		_ = pw.Stop()
		return err
	}
	a.pw = pw
	a.browser = browser
	a.ready = true
	return nil
}

func (a *PlaywrightAdapter) buildLaunchOptions() playwright.BrowserTypeLaunchOptions {
	opts := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(a.cfg.Headless),
		Args:     append([]string{}, a.cfg.ExtraArgs...),
	}

	if a.cfg.Channel != "" {
		opts.Channel = playwright.String(a.cfg.Channel)
	} else if a.cfg.ChromeChannel != "" {
		opts.Channel = playwright.String(a.cfg.ChromeChannel)
	}

	if a.cfg.DownloadsPath != "" {
		opts.DownloadsPath = playwright.String(a.cfg.DownloadsPath)
	}

	if a.cfg.ProxyConfig != nil && a.cfg.ProxyConfig.Server != "" {
		proxy := &playwright.Proxy{
			Server: a.cfg.ProxyConfig.Server,
		}
		if a.cfg.ProxyConfig.Username != "" {
			proxy.Username = playwright.String(a.cfg.ProxyConfig.Username)
		}
		if a.cfg.ProxyConfig.Password != "" {
			proxy.Password = playwright.String(a.cfg.ProxyConfig.Password)
		}
		opts.Proxy = proxy
	} else if a.cfg.Proxy != "" {
		opts.Proxy = &playwright.Proxy{Server: a.cfg.Proxy}
	}

	// Keep launch timeout aligned with crawler defaults until a dedicated browser launch timeout exists.
	opts.Timeout = playwright.Float(float64(config.DefaultPageTimeoutMs))

	// Expose Chromium's remote debugging port when configured.
	if a.cfg.DebuggingPort > 0 && (a.cfg.BrowserType == "" || a.cfg.BrowserType == "chromium") {
		opts.Args = append(opts.Args, fmt.Sprintf("--remote-debugging-port=%d", a.cfg.DebuggingPort))
	}

	return opts
}

func (a *PlaywrightAdapter) FetchHTML(ctx context.Context, url string, cfg config.CrawlerRunConfig) (FetchResult, error) {
	a.mu.Lock()
	if !a.ready || a.browser == nil {
		a.mu.Unlock()
		return FetchResult{}, stderrors.ErrBrowserNotStarted
	}
	b := a.browser
	a.mu.Unlock()

	page, err := b.NewPage()
	if err != nil {
		return FetchResult{}, err
	}
	defer page.Close()

	waitUntil := cfg.WaitUntil
	if waitUntil == "" {
		waitUntil = config.DefaultWaitUntil
	}

	timeoutMs := cfg.PageTimeoutMs
	if timeoutMs <= 0 {
		timeoutMs = config.DefaultPageTimeoutMs
	}
	timeout := float64(timeoutMs)
	waitUntilState := playwright.WaitUntilState(waitUntil)
	resp, err := page.Goto(
		url,
		playwright.PageGotoOptions{
			WaitUntil: &waitUntilState,
			Timeout:   &timeout,
		},
	)
	if err != nil {
		return FetchResult{}, err
	}

	if cfg.WaitFor != "" {
		waitForTimeout := timeout
		if cfg.WaitForTimeoutMs > 0 {
			waitForTimeout = float64(cfg.WaitForTimeoutMs)
		}
		if _, err := page.WaitForSelector(cfg.WaitFor, playwright.PageWaitForSelectorOptions{
			Timeout: &waitForTimeout,
		}); err != nil {
			return FetchResult{}, err
		}
	}

	// if context is already canceled return immediatly
	select {
	case <-ctx.Done():
		return FetchResult{}, ctx.Err()
	default:
	}

	html, err := page.Content()
	if err != nil {
		return FetchResult{}, err
	}

	result := FetchResult{
		HTML:          html,
		RedirectedURL: page.URL(),
	}
	if resp != nil {
		result.StatusCode = resp.Status()
		headers := make(map[string][]string, len(resp.Headers()))
		for k, v := range resp.Headers() {
			headers[k] = []string{v}
		}
		result.ResponseHeaders = headers
	}

	return result, nil
}

func (a *PlaywrightAdapter) Close(ctx context.Context) error {
	_ = ctx

	a.mu.Lock()
	if !a.ready {
		a.mu.Unlock()
		return nil
	}
	b := a.browser
	pw := a.pw
	a.browser = nil
	a.pw = nil
	a.ready = false
	a.mu.Unlock()

	if b != nil {
		if err := b.Close(); err != nil {
			return err
		}
	}
	if pw != nil {
		if err := pw.Stop(); err != nil {
			return err
		}
	}
	return nil
}
