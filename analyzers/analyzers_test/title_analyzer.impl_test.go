package analyzers_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/janithT/webpage-analyzer/analyzers"
)

func TestTitleAnalyzer_Simple(t *testing.T) {
	fmt.Println("dsdsdsdsds1")
	analyzer := analyzers.TitleAnalyzer()

	// Create HTML with a title
	html := `<html><head><title>Test Title</title></head><body></body></html>`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("failed to create goquery document: %v", err)
	}
	fmt.Println("dsdsdsdsds")
	result := analyzer.Analyze(doc, " ")

	// in above html check title key and value
	if result.Key != "title" {
		t.Errorf("expected key 'title', got %q", result.Key)
	}
	if result.Value != "Test Title" {
		t.Errorf("expected value 'Test Title', got %q", result.Value)
	}
}

func TestTitleAnalyzer_EmptyTitle(t *testing.T) {
	analyzer := analyzers.TitleAnalyzer()

	html := `<html><head><title></title></head><body></body></html>`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("failed to create goquery document: %v", err)
	}

	result := analyzer.Analyze(doc, "dsd")
	if result.Value != "" {
		t.Errorf("expected empty title, got %q", result.Value)
	}
}
