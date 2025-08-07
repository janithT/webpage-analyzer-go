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
		// Handle 403 Forbidden
		if status == http.StatusForbidden {
			responses.WriteError(ginC, http.StatusForbidden, "This URL is not accessible. Please check if the site is blocking requests.")
			return
		}

		// Handle DNS/host not found
		if strings.Contains(strings.ToLower(err.Error()), "No such host, Please enter correct URL") ||
			strings.Contains(strings.ToLower(err.Error()), "Server is misbehaving") ||
			strings.Contains(strings.ToLower(err.Error()), "lookup") {
			responses.WriteError(ginC, http.StatusBadRequest, "The specified host could not be found. Please check the domain name.")
			return
		}

		// Generic error
		responses.WriteError(ginC, status, err.Error())
		return
	}

	// Initializing the page analyzers.
	// Then pass to the pool for analyze each.
	analyzersList := []analyzers.Analyzer{
		analyzers.HTMLVersionAnalyzer(),
		analyzers.TitleAnalyzer(),
		analyzers.HeadingAnalyzer(),
		analyzers.LoginFormAnalyzer(),
		analyzers.LinkAnalyzer(),
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
