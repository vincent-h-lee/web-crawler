package crawler

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

func Crawl(u string, page rod.Page) (CrawlEvent, error) {
	timestamp := time.Now().UTC()
	log.Printf("Crawling: %s at %s", u, timestamp)
	timeoutDuration := 30 * time.Second

	/* e := proto.NetworkResponseReceived{}
	wait := page.Timeout(timeoutDuration).WaitEvent(&e) */

	err := page.Navigate(u)
	if err != nil {
		return CrawlEvent{}, errors.Join(errors.New("could not connect to page"), err)
	}

	err = page.Timeout(timeoutDuration).WaitStable(timeoutDuration)
	//wait()

	if err != nil {
		return CrawlEvent{}, errors.Join(errors.New("page could not stabilize"), err)
	}

	titleEl, err := page.Element("h1")
	if err != nil {
		return CrawlEvent{}, errors.Join(errors.New("page could not find h1 element"), err)
	}
	title, err := titleEl.Text()
	if err != nil {
		return CrawlEvent{}, errors.Join(errors.New("page could not find h1 text"), err)
	}

	var headings []Heading
	els, err := page.Elements("h1,h2,h3,h4,h5,h6")
	if err != nil {
		return CrawlEvent{}, errors.Join(errors.New("page could not find heading elements"), err)
	}
	for _, el := range els {
		headings = append(headings, Heading{Text: strings.TrimSpace(el.MustText()), Tag: el.MustProperty("tagName").String()})
	}

	var links []Link
	urlsSet := map[string]bool{}
	els, err = page.Elements("a")
	if err != nil {
		return CrawlEvent{}, errors.Join(errors.New("page could not find anchor elements"), err)
	}
	for _, el := range els {
		property, err := el.Property("href")
		if err != nil {
			// just skip if we can't find links
			continue
		}
		href := property.String()

		_, err = url.ParseRequestURI(href)

		if err == nil && !urlsSet[href] {
			urlsSet[href] = true
			links = append(links, Link{Url: href})
		}
	}

	// TODO store status code
	ev := NewCrawlEvent(u, title, headings, links, 0, timestamp)
	log.Printf("Finished crawling %s", u)
	return ev, nil
}
