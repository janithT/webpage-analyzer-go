package analyzers

import (
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const titleHtmlTag = "title"

type titleAnalyzer struct{}

// Construct function to title analyzer
func TitleAnalyzer() Analyzer {
	return &titleAnalyzer{}
}

// Logic here
func (a titleAnalyzer) Analyze(doc *goquery.Document, _ string) Result {

	startTime := time.Now()
	log.Println("Title analyzer started")
	defer func(start time.Time) {
		log.Printf("Title analyzer completed. Duration : %v ms", time.Since(start).Milliseconds())
	}(startTime)

	// Get the title of doc
	title := doc.Find(titleHtmlTag).Text()

	return Result{Key: "title", Value: title}
}
