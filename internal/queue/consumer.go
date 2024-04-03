package queue

import (
	"context"
	"log"
	"vincent-h-lee/web-crawler/internal/crawler"
)

func NewConsumer(repo crawler.CrawlerRepository, publisher Publisher, pool crawler.Pool) func(ctx context.Context, u string) error {
	return func(ctx context.Context, u string) error {
		log.Printf("Consuming url: %s", u)
		page := pool.Get()
		defer pool.Put(page)

		hasRecentlyCrawled, err := repo.HasRecentlyCrawled(ctx, u)
		log.Printf("Has recently crawled url: %t", hasRecentlyCrawled)
		if err != nil {
			return err
		}

		if hasRecentlyCrawled {
			// early exit don't need to crawl
			return nil
		}

		ev, err := crawler.Crawl(u, *page)
		if err != nil {
			return err
		}

		ev, err = repo.Store(ctx, ev)
		if err != nil {
			return err
		}

		for _, link := range ev.Links {
			log.Printf("Publishing url: %s", link.Url)
			err = publisher.Publish(ctx, link.Url)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
