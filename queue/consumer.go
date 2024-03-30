package queue

import (
	"log"
	"vincent-h-lee/web-crawler/crawler"
)

type consumer struct {
	repo *crawler.CrawlerRepository
}

func NewConsumer(repo *crawler.CrawlerRepository) *consumer {
	return &consumer{repo: repo}
}

func (c *consumer) Consume(u string) (crawler.CrawlEvent, error) {
	ev, err := crawler.Crawl(u)
	if err != nil {
		log.Fatal(err)
	}

	return c.repo.Store(ev)
}
