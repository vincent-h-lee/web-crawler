package queue

import (
	"context"
	"database/sql"
	"log"
	"vincent-h-lee/web-crawler/internal/crawler"
)

func NewConsumer(repo crawler.CrawlerRepository, publisher Publisher, pool crawler.Pool) func(ctx context.Context, u string) error {
	return func(ctx context.Context, u string) error {
		page := pool.Get()
		defer pool.Put(page)

		hasRecentlyCrawled, err := repo.HasRecentlyCrawled(ctx, u)
		if err != nil && err != sql.ErrNoRows {
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

		log.Printf("Found %d links", len(ev.Links))
		for _, link := range ev.Links {
			hasRecentlyCrawled, err := repo.HasRecentlyCrawled(ctx, link.Url)
			if err != nil {
				return err
			}
			if !hasRecentlyCrawled {
				err = publisher.Publish(ctx, link.Url)

				if err != nil {
					return err
				}
				log.Printf("Publishing url: %s", link.Url)
			}
		}

		return nil
	}
}
