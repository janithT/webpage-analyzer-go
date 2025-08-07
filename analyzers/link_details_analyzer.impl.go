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
	links   sync.Map
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

	l.links = sync.Map{}

	// Extract links
	l.prepare(doc, rawHTML)

	// Prepare workers
	wg.Add(l.getMapLength())

	// Push each link to channel
	l.links.Range(func(key, _ interface{}) bool {
		l.urlExec.Create().
			Build(
				key.(string),
				&wg,
				func(url string, status int, latency int64) {
					if val, ok := l.links.Load(url); ok {
						lp := val.(LinkProperty)
						lp.StatusCode = status
						lp.Latency = latency
						l.links.Store(url, lp)
					}
				},
			).
			PushChannel()
		return true
	})

	// Wait for all
	wg.Wait()

	// Collect results & counts
	var results []LinkProperty
	internalCount := 0
	externalCount := 0
	unknownCount := 0

	l.links.Range(func(_, value interface{}) bool {
		lp := value.(LinkProperty)
		results = append(results, lp)

		switch lp.Type {
		case Internal:
			internalCount++
		case External:
			externalCount++
		default:
			unknownCount++
		}
		return true
	})

	totalCount := internalCount + externalCount + unknownCount

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
		l.links.Store(url, LinkProperty{
			Url:  url,
			Type: getLinkType(url, baseDomain),
		})
	}
}

// number of links in the map
func (l *linkAnalyzer) getMapLength() int {
	count := 0
	l.links.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
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
	return strings.HasPrefix(trimmed, httpPref)
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
