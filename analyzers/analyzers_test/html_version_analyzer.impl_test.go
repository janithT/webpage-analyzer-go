package analyzers_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/janithT/webpage-analyzer/analyzers"
)

func TestHTMLVersionAnalyzer(t *testing.T) {
	tests := []struct {
		name         string
		rawHTML      string
		expectedVers string
	}{
		{
			name:         "HTML5 doctype",
			rawHTML:      "<!DOCTYPE html><html><head><title>Test</title></head><body></body></html>",
			expectedVers: "HTML5",
		},
		{
			name: "HTML 4.01 Transitional",
			rawHTML: `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" 
				"http://www.w3.org/TR/html4/loose.dtd">
				<html><head><title>Test</title></head><body></body></html>`,
			expectedVers: "HTML 4.01 Transitional",
		},
		{
			name:         "No doctype",
			rawHTML:      `<html><head><title>No Doctype</title></head><body></body></html>`,
			expectedVers: "Unknown",
		},
	}

	analyzer := analyzers.HTMLVersionAnalyzer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// doc is unused in analyzer logic, but required by method signature
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.rawHTML))
			if err != nil {
				t.Fatalf("failed to create goquery document: %v", err)
			}

			result := analyzer.Analyze(doc, tt.rawHTML)

			if result.Key != "htmlVersion" {
				t.Errorf("expected key 'htmlVersion', got %q", result.Key)
			}
			if result.Value != tt.expectedVers {
				t.Errorf("expected value %q, got %q", tt.expectedVers, result.Value)
			}
		})
	}
}
