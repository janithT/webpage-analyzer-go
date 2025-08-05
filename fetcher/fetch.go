package fetcher

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Regex url validator
var urlRegex = regexp.MustCompile(`^(https?:\/\/)?([a-zA-Z0-9\-]+\.)+[a-zA-Z]{2,}(:\d+)?(\/[^\s]*)?$`)

// Check url is valied
func IsValidURL(uri string) bool {
	parsed, err := url.ParseRequestURI(uri)
	return err == nil && parsed.Scheme != "" && parsed.Host != ""
}

// Check url regex is valied
func IsRegexValidURL(uri string) bool {
	url := strings.TrimSpace(uri)
	return urlRegex.MatchString(url)
}

// Fetch and parse the url
func FetchAndParse(uri string) (*goquery.Document, string, int, error) {

	parsedURL, err := url.Parse(uri)
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(uri)
	if err != nil {
		return nil, "", http.StatusBadGateway, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, "", resp.StatusCode, errors.New(resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", http.StatusInternalServerError, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, "", http.StatusInternalServerError, err
	}

	doc.Url = parsedURL

	return doc, string(bodyBytes), 200, nil
}
