package extract

import (
	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

func ToMarkdown(cleanHTML string, pageURL string) (string, error) {
	markdown, err := htmltomarkdown.ConvertString(cleanHTML)
	if err != nil {
		return "", err
	}
	return markdown, nil
}
