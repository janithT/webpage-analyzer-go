package models

// Title analyzer model
type TitleInfo struct {
	Text string `json:"text"`
}

// Heading tags model
type HeadingStat struct {
	TagName     string   `json:"tagName"`
	TagContents []string `json:"tagContents"`
	TagCount    int      `json:"tagCount"`
}

// Link details model
type LinkType int

const (
	Internal LinkType = iota
	External
)

type LinkProperty struct {
	Url        string
	Type       LinkType
	StatusCode int
	Latency    int64
}
