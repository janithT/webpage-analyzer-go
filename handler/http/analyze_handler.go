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

	// Validate URL
	if !fetcher.IsValidURL(url) || !fetcher.IsRegexValidURL(url) {
		responses.WriteError(ginC, http.StatusBadRequest, "Invalid URL format")
		return
	}

	doc, raw, status, err := fetcher.FetchAndParse(url)
	if err != nil {
		// Handle errors
		if status == http.StatusForbidden {
			responses.WriteError(ginC, http.StatusForbidden, "URL not accessible or blocked.")
			return
		}
		if strings.Contains(strings.ToLower(err.Error()), "no such host") ||
			strings.Contains(strings.ToLower(err.Error()), "server is misbehaving") ||
			strings.Contains(strings.ToLower(err.Error()), "lookup") {
			responses.WriteError(ginC, http.StatusBadRequest, "Host not found. Check domain name.")
			return
		}
		responses.WriteError(ginC, status, err.Error())
		return
	}

	analyzersList := []analyzers.Analyzer{
		analyzers.HTMLVersionAnalyzer(),
		analyzers.TitleAnalyzer(),
		analyzers.HeadingAnalyzer(),
		analyzers.LoginFormAnalyzer(),
		analyzers.LinkAnalyzer(),
	}

	results := pool.ExecuteAnalyzers(analyzersList, doc, raw)

	data := make(map[string]interface{})
	for _, result := range results {
		if result.Error != "" {
			data[result.Key] = map[string]interface{}{"error": result.Error}
			continue
		}
		data[result.Key] = result.Value
	}

	responses.WriteSuccess(ginC, "Analyzed successfully", data)
}
