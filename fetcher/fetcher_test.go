package fetcher

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// test the fetch and parse is valid
func TestFetchAndParseValid(t *testing.T) {

	//Test on sample html
	html := `<html><title>Hello</title><body><h1>Test</h1></body></html>`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, html)
	}))
	defer ts.Close()

	doc, raw, status, err := FetchAndParse(ts.URL)
	if err != nil || status != 200 {
		t.Fatalf("Expected success, got status %d err %v", status, err)
	}
	if !strings.Contains(raw, "<title>") {
		t.Errorf("Raw HTML not returned")
	}
	if doc.Find("h1").Text() != "Test" {
		t.Errorf("Parsed HTML missing h1")
	}
}

// test the fetch and parse with bad status
func TestFetchAndParseBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	}))
	defer ts.Close()

	_, _, status, err := FetchAndParse(ts.URL)
	if status != http.StatusNotFound || err == nil {
		t.Errorf("Expected NotFound error, got %d err %v", status, err)
	}
}
