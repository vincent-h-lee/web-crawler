package queue

import (
	"context"
	"vincent-h-lee/web-crawler/internal/crawler"

	"github.com/go-rod/rod"
)

func NewConsumer(repo crawler.CrawlerRepository, publisher Publisher, pool *rod.PagePool, create func() *rod.Page) func(ctx context.Context, u string) error {
	return func(ctx context.Context, u string) error {
		page := pool.Get(create)
		defer pool.Put(page)

		hasRecentlyCrawled, err := repo.HasRecentlyCrawled(ctx, u)
		if err != nil {
			return err
		}

		if hasRecentlyCrawled {
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
			err = publisher.Publish(ctx, link.Url)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
