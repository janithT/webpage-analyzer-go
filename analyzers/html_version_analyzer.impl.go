package analyzers

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type HTMLVersionAnalyzer struct{}

const defaultVersion string = "Unknown"

// htmlDoctypes maps known DOCTYPE identifiers to HTML version names.
var htmlDoctypes = map[string]string{
	"<!DOCTYPE html>":                        "HTML5",
	"-//W3C//DTD HTML 4.01//EN":              "HTML 4.01 Strict",
	"-//W3C//DTD HTML 4.01 Transitional//EN": "HTML 4.01 Transitional",
	"-//W3C//DTD HTML 4.01 Frameset//EN":     "HTML 4.01 Frameset",
	"-//W3C//DTD XHTML 1.0 Strict//EN":       "XHTML 1.0 Strict",
	"-//W3C//DTD XHTML 1.0 Transitional//EN": "XHTML 1.0 Transitional",
	"-//W3C//DTD XHTML 1.0 Frameset//EN":     "XHTML 1.0 Frameset",
	"-//W3C//DTD XHTML 1.1//EN":              "XHTML 1.1",
}

func (a HTMLVersionAnalyzer) Analyze(_ *goquery.Document, raw string) Result {

	startTime := time.Now()
	log.Println("Html version analyzer started")
	defer func(start time.Time) {
		log.Printf("Html version analyzer completed. Duration : %v ms", time.Since(start).Milliseconds())
	}(startTime)

	version := defaultVersion

	lowerRaw := strings.ToLower(raw)

	// Find the correct version name
	for doctype, name := range htmlDoctypes {
		if strings.Contains(lowerRaw, strings.ToLower(doctype)) {
			version = name
			break
		}
	}
	return Result{Key: "htmlVersion", Value: version}
}
