package crawler

import (
	"errors"
	"net/url"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type link struct {
	Url string
}

type heading struct {
	Text string
	Tag  string
}

type DocumentData struct {
	// Description
	// Images?
	Title    string
	Headings []heading
	Links    []link
}

func Crawl(u string) (*DocumentData, error) {
	page, err := rod.New().MustConnect().Page(proto.TargetCreateTarget{URL: u})
	if err != nil {
		return nil, errors.Join(errors.New("could not connect to page"), err)
	}
	page.MustWaitStable()

	title := page.MustElement("h1").MustText()

	var headings []heading
	els := page.MustElements("h1,h2,h3,h4,h5,h6")
	for _, el := range els {
		headings = append(headings, heading{Text: strings.TrimSpace(el.MustText()), Tag: el.MustProperty("tagName").String()})
	}

	var links []link
	urlsSet := map[string]bool{}
	els = page.MustElements("a")
	for _, el := range els {
		href := el.MustProperty("href").String()
		_, err := url.ParseRequestURI(href)

		if err == nil && !urlsSet[u] {
			urlsSet[href] = true
			links = append(links, link{Url: el.MustProperty("href").String()})
		}
	}

	doc := &DocumentData{Title: title, Headings: headings, Links: links}
	return doc, nil
}
