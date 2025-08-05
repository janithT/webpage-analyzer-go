package http

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/janithT/webpage-analyzer/analyzers"
	"github.com/janithT/webpage-analyzer/fetcher"
	"github.com/janithT/webpage-analyzer/pool"
	"github.com/janithT/webpage-analyzer/responses"
)

// Main analyze handler
func AnalyzeHandler(ginC *gin.Context) {

	// Get url parameter
	url := strings.TrimSpace(ginC.Query("url"))
	log.Printf("Trimmed url = %v", url)
	// Check url is valied
	if !fetcher.IsValidURL(url) {
		responses.WriteError(ginC, http.StatusBadRequest, "Invalid URL format")
		return
	}
	// Check url is valied - regEx
	if !fetcher.IsRegexValidURL(url) {
		responses.WriteError(ginC, http.StatusBadRequest, "URL is not valied, Please check your URL.")
		return
	}
	// For fetch errors when fetching and parsing.
	doc, raw, status, err := fetcher.FetchAndParse(url)
	if err != nil {
		responses.WriteError(ginC, status, err.Error())
		return
	}

	// Initializing the page analyzers.
	// Then pass to the pool for analyze each.
	analyzersList := []analyzers.Analyzer{
		analyzers.HTMLVersionAnalyzer{},
		analyzers.TitleAnalyzer{},
		analyzers.HeadingAnalyzer{},
		analyzers.LoginFormAnalyzer{},
		analyzers.LinkAnalyzer{},
	}

	// Add pool to execute
	results := pool.ExecuteAnalyzers(analyzersList, doc, raw)

	// Filling the data object
	data := make(map[string]interface{})
	for _, result := range results {
		if result.Error != "" {
			data[result.Key] = map[string]interface{}{
				"error": result.Error,
			}
			continue
		}
		data[result.Key] = result.Value
	}

	responses.WriteSuccess(ginC, "Analyzed successfully", data)
}
