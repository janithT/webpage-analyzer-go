package analyzers

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const (
	inputTag          = "input"
	typeAttribute     = "type"
	passwordAttribute = "password"
)

type LoginFormAnalyzer struct{}

func (a LoginFormAnalyzer) Analyze(doc *goquery.Document, rawHTML string) Result {

	startTime := time.Now()
	log.Println("Login exists analyzer started")
	defer func(start time.Time) {
		log.Printf("Login exists analyzer completed. Duration : %v ms", time.Since(start).Milliseconds())
	}(startTime)

	tokenizer := html.NewTokenizer(strings.NewReader(rawHTML))

	for {
		switch tokenizer.Next() {
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == inputTag {
				if attrVal := getAttr(token, typeAttribute); strings.ToLower(attrVal) == passwordAttribute {
					// Found password field, return immediately
					return Result{Key: "hasLoginForm", Value: true}
				}
			}
		case html.ErrorToken:
			// End of document or error
			return Result{Key: "hasLoginForm", Value: false}
		}
	}
}

func getAttr(token html.Token, attrName string) string {
	for _, attr := range token.Attr {
		if strings.EqualFold(attr.Key, attrName) {
			return attr.Val
		}
	}
	return ""
}
