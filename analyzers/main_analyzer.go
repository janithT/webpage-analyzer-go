package analyzers

import "github.com/PuerkitoBio/goquery"

type Result struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Error string      `json:"error,omitempty"`
}

// main interface for analyzers
type Analyzer interface {
	Analyze(doc *goquery.Document, raw string) Result
}
