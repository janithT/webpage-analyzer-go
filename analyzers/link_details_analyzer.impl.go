package analyzers

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

// LinkType > internal or external
type LinkType int

const (
	Internal LinkType = iota
	External
)

// Get string form for JSON output
func (lt LinkType) MarshalJSON() ([]byte, error) {
	switch lt {
	case Internal:
		return []byte(`"Internal"`), nil
	case External:
		return []byte(`"External"`), nil
	default:
		return []byte(`"Unknown"`), nil
	}
}

// LinkProperty stores each link's info
type LinkProperty struct {
	Url        string   `json:"url"`
	Type       LinkType `json:"type"`
	StatusCode int      `json:"status_code"`
	Latency    int64    `json:"latency"`
}

// LinkAnalyzer is the analyzer struct
type LinkAnalyzer struct{}

// Analyze performs full analysis of the document and HTML
func (a LinkAnalyzer) Analyze(doc *goquery.Document, rawHTML string) Result {
	startTime := time.Now()
	log.Println("Link analyzer started")

	var links sync.Map
	baseHost := ""
	baseDomain := ""
	log.Printf("basedoc.Urlost = %v", doc)
	if doc.Url != nil {
		baseHost = doc.Url.Host
		log.Printf("baseHost = %v", baseHost)
		if baseHost != "" {
			if domain, err := publicsuffix.EffectiveTLDPlusOne(baseHost); err == nil {
				baseDomain = domain
			}
		}
	}

	tokenizer := html.NewTokenizer(strings.NewReader(rawHTML))
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.StartTagToken {
			token := tokenizer.Token()
			switch token.Data {
			case "a":
				a.processToken(&links, token, "href", baseDomain)
			case "script":
				a.processToken(&links, token, "src", baseDomain)
			case "link":
				a.processToken(&links, token, "href", baseDomain)
			}
		}
	}

	var wg sync.WaitGroup
	links.Range(func(key, value interface{}) bool {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			status, latency := checkLink(link)
			if val, ok := links.Load(link); ok {
				lp := val.(LinkProperty)
				lp.StatusCode = status
				lp.Latency = latency
				links.Store(link, lp)
			}
		}(key.(string))
		return true
	})

	wg.Wait()
	log.Printf("Link analyzer completed. Duration: %v ms", time.Since(startTime).Milliseconds())

	var results []LinkProperty
	links.Range(func(_, value interface{}) bool {
		results = append(results, value.(LinkProperty))
		return true
	})

	return Result{
		Key:   "urls",
		Value: results,
	}
}

// processToken handles individual tag attributes and normalizes URL
func (a LinkAnalyzer) processToken(links *sync.Map, token html.Token, attrName, baseDomain string) {
	attrVal := getLinkAttr(token, attrName)
	if attrVal == "" {
		return
	}

	u, err := url.Parse(attrVal)
	if err != nil {
		return
	}

	var absURL string
	if u.IsAbs() {
		absURL = u.String()
	} else {
		// Skip relative URLs for now
		return
	}

	// Skip invalid or non-http links
	if !strings.HasPrefix(absURL, "http") {
		return
	}

	linkHost := getHost(absURL)
	linkDomain, err := publicsuffix.EffectiveTLDPlusOne(linkHost)
	log.Printf("link host - %v | domain - %v | base domain - %v", linkHost, linkDomain, baseDomain)

	linkType := External
	if err == nil && baseDomain != "" && linkDomain == baseDomain {
		linkType = Internal
	}

	links.Store(absURL, LinkProperty{
		Url:  absURL,
		Type: linkType,
	})
}

// getHost parses a URL and returns the host
func getHost(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return u.Host
}

// getLinkAttr retrieves attribute from token
func getLinkAttr(token html.Token, attrName string) string {
	for _, attr := range token.Attr {
		if strings.EqualFold(attr.Key, attrName) {
			return attr.Val
		}
	}
	return ""
}

// checkLink makes an HTTP GET request and returns status + latency
func checkLink(link string) (int, int64) {
	start := time.Now()
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(link)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return 0, latency
	}
	defer resp.Body.Close()
	return resp.StatusCode, latency
}
