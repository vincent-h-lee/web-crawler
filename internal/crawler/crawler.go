package crawler

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func Crawl(u string, page rod.Page) (CrawlEvent, error) {
	timestamp := time.Now().UTC()
	log.Printf("Crawling: %s at %s", u, timestamp)

	e := proto.NetworkResponseReceived{}
	wait := page.WaitEvent(&e)

	err := page.Navigate(u)
	if err != nil {
		return CrawlEvent{}, errors.Join(errors.New("could not connect to page"), err)
	}
	page.MustWaitStable()
	wait()

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

		if err == nil && !urlsSet[href] {
			urlsSet[href] = true
			links = append(links, Link{Url: el.MustProperty("href").String()})
		}
	}

	ev := NewCrawlEvent(u, title, headings, links, e.Response.Status, timestamp)
	log.Printf("Finished crawling %s", u)
	return ev, nil
}
