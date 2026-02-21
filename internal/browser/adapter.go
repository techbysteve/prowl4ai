package browser

import (
	"context"

	"github.com/techbysteve/prowl4ai/internal/config"
)

type FetchResult struct {
	HTML            string
	StatusCode      int
	RedirectedURL   string
	ResponseHeaders map[string][]string
}

type Adapter interface {
	Start(ctx context.Context) error
	FetchHTML(ctx context.Context, url string, cfg config.CrawlerRunConfig) (FetchResult, error)
	Close(ctx context.Context) error
}
