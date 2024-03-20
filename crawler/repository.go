package crawler

import "time"

type crawlEvent struct {
	Url  string
	Data DocumentData
	Date time.Time
}

func NewCrawlEvent(u string, data DocumentData, now time.Time) crawlEvent {
	return crawlEvent{Url: u, Data: data, Date: now}
}

type CrawlerRepository struct{}

func (r *CrawlerRepository) Store(ev crawlEvent) (*crawlEvent, error) {
	// TODO
	return nil, nil
}

func (r *CrawlerRepository) Get(id string) (*crawlEvent, error) {
	// TODO
	return nil, nil
}
