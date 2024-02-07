package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

func ParseURL(rawURL string) string {
	urll, err := parseURL(rawURL)
	if err != nil {
		log.Error().
			Err(err).
			Str("url", urll).
			Msg("parse.ParseURL(): couldn't parse url")
		return rawURL
	}
	return urll
}

func parseURL(rawURL string) (string, error) {
	// rawURL may be empty string, function should return empty string then.
	rawURL = strings.TrimSpace(rawURL)
	parsedURL, parseErr := url.Parse(rawURL)
	if parseErr != nil {
		return "", fmt.Errorf("parse.parseURL(): failed url.Parse() on url(%v). error: %w", rawURL, parseErr)
	}

	urlString := parsedURL.String()
	if len(urlString) != 0 && len(parsedURL.Path) == 0 { // https://example.org -> https://example.org/
		urlString += "/"
	}

	return urlString, nil
}

func ParseTextWithHTML(rawHTML string) string {
	text, err := parseTextWithHTML(rawHTML)
	if err != nil {
		log.Error().
			Err(err).
			Str("html", rawHTML).
			Msg("parse.ParseTextWithHTML(): failed parsing text with html")
		return rawHTML
	}
	return text
}

func parseTextWithHTML(rawHTML string) (string, error) {
	var result string = ""
	htmlNode, perr := html.ParseFragment(strings.NewReader(rawHTML), nil)
	if perr != nil {
		return "", fmt.Errorf("parse.parseTextWithHTML(): failed html.ParseFragment on %v. error: %w", rawHTML, perr)
	}
	for _, el := range htmlNode {
		sel := goquery.NewDocumentFromNode(el)
		result += sel.Text()
	}
	return result, nil
}
