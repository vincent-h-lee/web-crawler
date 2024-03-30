package crawler

import (
	"time"

	"github.com/uptrace/bun"
)

type CrawlEvent struct {
	bun.BaseModel `bun:"table:crawls"`

	Id         int64     `bun:",pk,autoincrement" json:"id"`
	Url        string    `json:"url"`
	StatusCode int       `json:"statusCode"`
	Timestamp  time.Time `json:"timestamp"`
	Title      string    `json:"title"`

	Headings []Heading `json:"headings" bun:"rel:has-many,join:id=crawl_id"`
	Links    []Link    `json:"links" bun:"rel:has-many,join:id=crawl_id"`
}

func NewCrawlEvent(u string, title string, headings []Heading, links []Link, statusCode int, timestamp time.Time) CrawlEvent {
	return CrawlEvent{
		Url:        u,
		StatusCode: statusCode,
		Timestamp:  timestamp,
		Title:      title,
		Headings:   headings,
		Links:      links,
	}
}
