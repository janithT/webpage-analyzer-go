package pool

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/janithT/webpage-analyzer/analyzers"
)

type MockAnalyzer struct{ key, val string }

func (m MockAnalyzer) Analyze(doc *goquery.Document, raw string) analyzers.Result {
	return analyzers.Result{Key: m.key, Value: m.val}
}

func (m MockAnalyzer) Key() string { return m.key }

func TestExecuteAnalyzers(t *testing.T) {
	anList := []analyzers.Analyzer{
		MockAnalyzer{"one", "val1"},
		MockAnalyzer{"two", "val2"},
	}

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader("<html/>"))
	results := ExecuteAnalyzers(anList, doc, "")
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	expected := map[string]string{"one": "val1", "two": "val2"}
	for _, r := range results {
		if expected[r.Key] != r.Value {
			t.Errorf("Expected %s -> %s, got %v", r.Key, expected[r.Key], r.Value)
		}
	}
}
