package analyzers

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// LinkType indicates internal or external
type LinkType int

const (
	Internal LinkType = iota
	External
	Unknown
)

// MarshalJSON for pretty printing LinkType
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

type LinkProperty struct {
	Url        string   `json:"url"`
	Type       LinkType `json:"type"`
	StatusCode int      `json:"status_code"`
	Latency    int64    `json:"latency"` // milliseconds
}

type linkAnalyzer struct {
	links []LinkProperty
	mu    sync.Mutex
}

// NewLinkAnalyzer creates a new linkAnalyzer instance
func LinkAnalyzer() Analyzer {
	return &linkAnalyzer{}
}

// Analyze extracts all URLs and fetches their status asynchronously
func (l *linkAnalyzer) Analyze(doc *goquery.Document, rawHTML string) Result {
	startTime := time.Now()
	log.Println("Link analyzer started")
	defer func(start time.Time) {
		log.Printf("Link analyzer completed. Duration: %v ms", time.Since(start).Milliseconds())
	}(startTime)

	l.links = nil

	baseDomain := ""
	if doc.Url != nil && doc.Url.Host != "" {
		baseDomain = doc.Url.Hostname()
	}

	// Temporary map to deduplicate URLs
	linkMap := make(map[string]LinkProperty)

	// Find all tags with URLs (a[href], link[href], script[src], img[src])
	doc.Find("a[href], link[href], script[src], img[src]").Each(func(i int, s *goquery.Selection) {
		var attr string
		if s.Is("a, link") {
			attr, _ = s.Attr("href")
		} else { // script, img
			attr, _ = s.Attr("src")
		}
		if attr == "" {
			return
		}

		absUrl := resolveUrl(doc.Url, attr)
		if absUrl == "" {
			return
		}

		// Check if URL is valid http or https
		parsed, err := url.Parse(absUrl)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			return
		}

		// Deduplicate
		if _, exists := linkMap[absUrl]; !exists {
			linkMap[absUrl] = LinkProperty{
				Url:  absUrl,
				Type: getLinkType(absUrl, baseDomain),
			}
		}
	})

	// Convert map to slice
	for _, lp := range linkMap {
		l.links = append(l.links, lp)
	}

	// Use WaitGroup to track async checking of URLs
	var wg sync.WaitGroup
	wg.Add(len(l.links))

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for i := range l.links {
		go func(idx int) {
			defer wg.Done()
			start := time.Now()
			resp, err := client.Head(l.links[idx].Url)
			if err != nil || resp == nil {
				// Possibly try GET if HEAD fails:
				respGet, errGet := client.Get(l.links[idx].Url)
				if errGet != nil || respGet == nil {
					// Mark as unreachable; status code 0
					l.mu.Lock()
					l.links[idx].StatusCode = 0
					l.links[idx].Latency = int64(time.Since(start).Milliseconds())
					l.mu.Unlock()
					return
				}
				resp = respGet
				defer resp.Body.Close()
			} else {
				defer resp.Body.Close()
			}
			latency := int64(time.Since(start).Milliseconds())

			l.mu.Lock()
			l.links[idx].StatusCode = resp.StatusCode
			l.links[idx].Latency = latency
			l.mu.Unlock()
		}(i)
	}

	wg.Wait()

	// Counts
	internalCount := 0
	externalCount := 0
	unknownCount := 0

	for _, link := range l.links {
		switch link.Type {
		case Internal:
			internalCount++
		case External:
			externalCount++
		default:
			unknownCount++
		}
	}

	totalCount := len(l.links)

	return Result{
		Key: "urls",
		Value: map[string]interface{}{
			"total_count":    totalCount,
			"internal_count": internalCount,
			"external_count": externalCount,
			"unknown_count":  unknownCount,
			"links":          l.links,
		},
	}
}

// resolveUrl resolves a possibly relative URL href against the base *url.URL
func resolveUrl(base *url.URL, href string) string {
	if base == nil {
		return href
	}
	u, err := base.Parse(href)
	if err != nil {
		return ""
	}
	return u.String()
}

// getLinkType determines if URL is Internal, External or Unknown given baseDomain
func getLinkType(link string, baseDomain string) LinkType {
	parsedUrl, err := url.Parse(link)
	if err != nil || baseDomain == "" {
		return Unknown
	}
	host := parsedUrl.Hostname()
	if host == "" {
		return Unknown
	}
	if host == baseDomain {
		return Internal
	}
	return External
}
