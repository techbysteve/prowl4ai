package prowler

import (
	"context"
	"sync"

	"github.com/techbysteve/prowl4ai/internal/browser"
	"github.com/techbysteve/prowl4ai/internal/config"
	"github.com/techbysteve/prowl4ai/internal/extract"
	"github.com/techbysteve/prowl4ai/internal/model"
	"github.com/techbysteve/prowl4ai/internal/stderrors"
)

type Service struct {
	browser browser.Adapter
	mu      sync.Mutex
	ready   bool
}

func NewService(adapter browser.Adapter) *Service {
	return &Service{
		browser: adapter,
	}
}

func (s *Service) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ready {
		return nil
	}
	if s.browser == nil {
		return stderrors.ErrBrowserNotStarted
	}
	if err := s.browser.Start(ctx); err != nil {
		return err
	}
	s.ready = true
	return nil
}

func (s *Service) FetchHTML(ctx context.Context, url string, cfg config.CrawlerRunConfig) (browser.FetchResult, error) {
	if url == "" {
		return browser.FetchResult{}, stderrors.ErrInvalidURL
	}

	if err := s.Start(ctx); err != nil {
		return browser.FetchResult{}, err
	}

	s.mu.Lock()
	if !s.ready {
		s.mu.Unlock()
		return browser.FetchResult{}, stderrors.ErrServiceNotReady
	}
	if s.browser == nil {
		s.mu.Unlock()
		return browser.FetchResult{}, stderrors.ErrBrowserNotStarted
	}
	b := s.browser
	s.mu.Unlock()

	if cfg.PageTimeoutMs <= 0 {
		cfg.PageTimeoutMs = config.DefaultPageTimeoutMs
	}
	if cfg.WaitUntil == "" {
		cfg.WaitUntil = config.DefaultWaitUntil
	}

	result, err := b.FetchHTML(ctx, url, cfg)
	if err != nil {
		return browser.FetchResult{}, err
	}

	return result, nil
}

func (s *Service) Close(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.ready {
		return nil
	}
	if s.browser == nil {
		s.ready = false
		return nil
	}
	if err := s.browser.Close(ctx); err != nil {
		return err
	}
	s.ready = false
	return nil
}

func (s *Service) Run(ctx context.Context, url string, cfg config.CrawlerRunConfig) (model.CrawlResult, error) {
	fetchResult, err := s.FetchHTML(ctx, url, cfg)
	if err != nil {
		return model.CrawlResult{
			URL:          url,
			Success:      false,
			ErrorMessage: err.Error(),
		}, err
	}

	headers := make(map[string]any, len(fetchResult.ResponseHeaders))
	for k, v := range fetchResult.ResponseHeaders {
		headers[k] = v
	}

	baseURL := fetchResult.RedirectedURL
	if baseURL == "" {
		baseURL = url
	}

	extractOut, extractErr := extract.Process(fetchResult.HTML, baseURL, cfg)
	if extractErr != nil {
		return model.CrawlResult{
			URL:             url,
			HTML:            fetchResult.HTML,
			Success:         false,
			ErrorMessage:    extractErr.Error(),
			ResponseHeaders: headers,
			StatusCode:      fetchResult.StatusCode,
			RedirectedURL:   fetchResult.RedirectedURL,
		}, extractErr
	}

	return model.CrawlResult{
		URL:             url,
		HTML:            fetchResult.HTML,
		CleanedHTML:     extractOut.CleanedHTML,
		Success:         true,
		Markdown:        model.Markdown(extractOut.Markdown),
		Metadata:        extractOut.Metadata,
		ResponseHeaders: headers,
		StatusCode:      fetchResult.StatusCode,
		RedirectedURL:   fetchResult.RedirectedURL,
	}, nil
}
