package analyzers

import (
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/janithT/webpage-analyzer/handler/models"
)

type HeadingAnalyzer struct{}

func NewHeadingAnalyzer() Analyzer {
	return &HeadingAnalyzer{}
}

func (a HeadingAnalyzer) Analyze(doc *goquery.Document, _ string) Result {

	startTime := time.Now()
	log.Println("Heading analyzer started")
	defer func(start time.Time) {
		log.Printf("Heading analyzer completed. Duration : %v ms", time.Since(start).Milliseconds())
	}(startTime)

	var headingStats []models.HeadingStat

	// I'm using GoQuery for this Google suggests it's best tools for tag analysis.
	for i := 1; i <= 6; i++ {
		tag := "h" + string('0'+i)
		var tagContents []string
		doc.Find(tag).Each(func(_ int, s *goquery.Selection) {
			text := s.Text()
			if text != "" {
				tagContents = append(tagContents, text)
			}
		})

		if len(tagContents) > 0 {
			headingStats = append(headingStats, models.HeadingStat{
				TagName:     tag,
				TagContents: tagContents,
				TagCount:    len(tagContents),
			})
		}
	}

	return Result{
		Key:   "headings",
		Value: headingStats,
	}
}
