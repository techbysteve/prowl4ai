package model

type Markdown string

type CrawlResult struct {
	URL             string         `json:"url"`
	HTML            string         `json:"html,omitempty"`
	CleanedHTML     string         `json:"cleaned_html,omitempty"`
	Success         bool           `json:"success"`
	Markdown        Markdown       `json:"markdown,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
	ErrorMessage    string         `json:"error_message,omitempty"`
	SessionID       string         `json:"session_id,omitempty"`
	ResponseHeaders map[string]any `json:"response_headers,omitempty"`
	StatusCode      int            `json:"status_code,omitempty"`
	RedirectedURL   string         `json:"redirected_url,omitempty"`
}
