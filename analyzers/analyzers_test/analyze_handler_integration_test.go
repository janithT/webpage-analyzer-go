package analyzers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	myhttp "github.com/janithT/webpage-analyzer/handler/http"
)

var resp struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func TestAnalyzeHandler_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test router
	router := gin.New()
	router.GET("/analyze", myhttp.AnalyzeHandler)

	// Serve a simple HTML page via a test HTTP server
	testHTML := "https://example.com"

	fmt.Printf(" url test = %v", url.QueryEscape(testHTML))
	// Create request to /analyze with our test server's URL
	reqURL := "/analyze?url=" + url.QueryEscape(testHTML)

	req, _ := http.NewRequest("GET", reqURL, nil)

	// Record the response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verify HTTP status
	// if w.Code != http.StatusOK {
	// 	t.Fatalf("Expected status 200, got %d", w.Code)
	// }

	// Verify response body contains both title and htmlVersion
	body := w.Body.String()

	fmt.Printf(" body response = %v", body)
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	// Check the title of test url
	if title, ok := resp.Data["title"].(string); !ok || title != "Example Domain" {
		t.Errorf("Expected title 'Example Domain', got %v", resp.Data["title"])
	}

	// check the version == HTML5
	if version, ok := resp.Data["htmlVersion"].(string); !ok || version != "HTML5" {
		t.Errorf("Expected htmlVersion 'HTML5', got %v", resp.Data["htmlVersion"])
	}
}
