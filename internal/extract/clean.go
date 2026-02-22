package extract

import (
	"bytes"
	"net/url"
	"strings"

	readability "codeberg.org/readeck/go-readability/v2"
)

func CleanHTML(rawHTML string, pageURL string, onlyText bool) (string, map[string]any, error) {
	parsedUrl, err := url.ParseRequestURI(pageURL)
	if err != nil {
		return "", nil, err
	}
	readableContent, err := readability.FromReader(strings.NewReader(rawHTML), parsedUrl)
	if err != nil {
		return "", nil, err
	}

	var htmlBuffer bytes.Buffer

	if err := readableContent.RenderHTML(&htmlBuffer); err != nil {
		return "", nil, err
	}

	metadata := map[string]any{
		"title":   readableContent.Title(),
		"byline":  readableContent.Byline(),
		"excerpt": readableContent.Excerpt(),
		"lang":    readableContent.Language(),
	}

	return htmlBuffer.String(), metadata, nil
}
