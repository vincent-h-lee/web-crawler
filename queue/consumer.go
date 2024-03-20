package queue

import (
	"log"
	"time"
	"vincent-h-lee/web-crawler/crawler"
)

// TODO browser pool https://github.com/go-rod/rod/blob/46baf3aad803ed5cd8671aa325cbae4e297a89a4/setup_test.go#L59
func Consume(u string) {
	now := time.Now()
	doc, err := crawler.Crawl(u)
	if err != nil {
		log.Fatal(err)
	}

	repo := crawler.CrawlerRepository{}
	repo.Store(crawler.NewCrawlEvent(u, *doc, now))
}
