package crawler

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/uptrace/bun"
)

type Link struct {
	bun.BaseModel `bun:"table:crawl_links"`
	Url           string `json:"url"`
	Text          string `json:"text"`

	CrawlId int64 `json:"-"`
}

type Heading struct {
	bun.BaseModel `bun:"table:crawl_headings"`

	Text    string `json:"text"`
	Tag     string `json:"tag"`
	CrawlId int64  `json:"-"`
}

func Crawl(u string) (CrawlEvent, error) {
	timestamp := time.Now().UTC()
	page, err := rod.New().MustConnect().Page(proto.TargetCreateTarget{URL: u})
	if err != nil {
		return CrawlEvent{}, errors.Join(errors.New("could not connect to page"), err)
	}
	page.MustWaitStable()

	title := page.MustElement("h1").MustText()

	var headings []Heading
	els := page.MustElements("h1,h2,h3,h4,h5,h6")
	for _, el := range els {
		headings = append(headings, Heading{Text: strings.TrimSpace(el.MustText()), Tag: el.MustProperty("tagName").String()})
	}

	var links []Link
	urlsSet := map[string]bool{}
	els = page.MustElements("a")
	for _, el := range els {
		href := el.MustProperty("href").String()
		_, err := url.ParseRequestURI(href)

		if err == nil && !urlsSet[u] {
			urlsSet[href] = true
			links = append(links, Link{Url: el.MustProperty("href").String()})
		}
	}

	ev := NewCrawlEvent(u, title, headings, links, 0, timestamp)
	return ev, nil
}
