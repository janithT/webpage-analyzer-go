package analyzers

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	channels "github.com/janithT/webpage-analyzer/channel"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

// LinkType indicates internal or external
type LinkType int

const (
	Internal LinkType = iota
	External
	Unknown
)

const (
	aTag      = "a"
	hrefAttr  = "href"
	httpPref  = "http"
	scriptTag = "script"
	srcAttr   = "src"
	linkTag   = "link"
)

type LinkProperty struct {
	Url        string   `json:"url"`
	Type       LinkType `json:"type"`
	StatusCode int      `json:"status_code"`
	Latency    int64    `json:"latency"`
}

func (lt LinkType) MarshalJSON() ([]byte, error) {
	switch lt {
	case Internal:
		return []byte(`"Internal"`), nil
	case External:
		return []byte(`"External"`), nil
	case Unknown:
		return []byte(`"Unknown"`), nil
	default:
		return []byte(`"Unknown"`), nil
	}
}

var wg sync.WaitGroup

type linkAnalyzer struct {
	links   []LinkProperty
	mu      sync.Mutex
	urlExec channels.UrlWorker
}

func LinkAnalyzer() Analyzer {
	return &linkAnalyzer{
		urlExec: channels.NewUrlWorker(),
	}
}

func (l *linkAnalyzer) Analyze(doc *goquery.Document, rawHTML string) Result {
	startTime := time.Now()

	log.Println("Link analyzer started")
	defer func(start time.Time) {
		log.Printf("Link analyzer completed. Duration: %v ms", time.Since(start).Milliseconds())
	}(startTime)

	l.links = []LinkProperty{}

	// Extract links
	l.prepare(doc, rawHTML)

	// Prepare workers
	wg := sync.WaitGroup{}
	wg.Add(len(l.links))

	// Push each link to channel
	for i := range l.links {
		idx := i
		url := l.links[idx].Url
		l.urlExec.Create().
			Build(
				url,
				&wg,
				func(_ string, status int, latency int64) {
					l.mu.Lock()
					l.links[idx].StatusCode = status
					l.links[idx].Latency = latency
					l.mu.Unlock()
				},
			).
			PushChannel()
	}

	// Wait for all
	wg.Wait()

	// Collect results & counts
	var results []LinkProperty
	internalCount := 0
	externalCount := 0
	unknownCount := 0

	for _, lp := range l.links {
		switch lp.Type {
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
			"links":          results,
		},
	}

	// // Collect results
	// var results []LinkProperty
	// l.links.Range(func(_, value interface{}) bool {
	// 	results = append(results, value.(LinkProperty))
	// 	return true
	// })

	// return Result{
	// 	Key:   "urls",
	// 	Value: results,
	// }
}

// prepare for link extraction
func (l *linkAnalyzer) prepare(doc *goquery.Document, rawHTML string) {
	baseDomain := ""
	if doc.Url != nil && doc.Url.Host != "" {
		if domain, err := publicsuffix.EffectiveTLDPlusOne(doc.Url.Host); err == nil {
			baseDomain = domain
		}
	}

	tokenizer := html.NewTokenizer(strings.NewReader(rawHTML))
	for {
		switch tokenizer.Next() {
		case html.StartTagToken:
			token := tokenizer.Token()
			switch token.Data {
			case aTag:
				l.storeIfValid(getTagAttribute(token, hrefAttr), baseDomain)
			case scriptTag:
				l.storeIfValid(getTagAttribute(token, srcAttr), baseDomain)
			case linkTag:
				l.storeIfValid(getTagAttribute(token, hrefAttr), baseDomain)
			}
		case html.ErrorToken:
			return
		}
	}
}

// if valied store the link
func (l *linkAnalyzer) storeIfValid(url, baseDomain string) {
	if isValidLink(url) {
		lp := LinkProperty{
			Url:  url,
			Type: getLinkType(url, baseDomain),
		}
		l.mu.Lock()
		l.links = append(l.links, lp)
		l.mu.Unlock()
	}
}

// number of links in the slice
func (l *linkAnalyzer) getMapLength() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.links)
}

// get link type > with base domain
func getLinkType(link, baseDomain string) LinkType {
	if baseDomain != "" && strings.Contains(link, baseDomain) {
		return Internal
	}
	return External
}

// check if the link is valid
func isValidLink(link string) bool {
	trimmed := strings.TrimSpace(link)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return false // skip empty and fragment-only links
	}
	if strings.HasPrefix(trimmed, "http") || strings.HasPrefix(trimmed, "/") || strings.HasPrefix(trimmed, "//") {
		return true
	}
	return false
}

// get tag attribute value
func getTagAttribute(token html.Token, attrName string) string {
	for _, attr := range token.Attr {
		if strings.EqualFold(attr.Key, attrName) {
			return attr.Val
		}
	}
	return ""
}
