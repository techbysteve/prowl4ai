package extract

import "github.com/techbysteve/prowl4ai/internal/config"

type Output struct {
	CleanedHTML string
	Markdown    string
	Metadata    map[string]any
}

func Process(rawHTML, baseURL string, cfg config.CrawlerRunConfig) (Output, error) {
	out := Output{
		CleanedHTML: rawHTML,
		Metadata:    map[string]any{},
	}

	if cfg.EnableCleanHTML {
		cleaned, metadata, err := CleanHTML(rawHTML, baseURL, cfg.OnlyText)
		if err != nil {
			return out, err
		}
		if cleaned != "" {
			out.CleanedHTML = cleaned
		}
		if metadata != nil {
			out.Metadata = metadata
		}
	}

	if cfg.EnableMarkdown {
		input := out.CleanedHTML
		if input == "" {
			input = rawHTML
		}

		markdown, err := ToMarkdown(input, baseURL)
		if err != nil {
			return out, err
		}
		out.Markdown = markdown
	}

	return out, nil
}
