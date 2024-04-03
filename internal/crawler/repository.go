package crawler

import (
	"context"
	"log"

	"github.com/uptrace/bun"
)

type CrawlerRepository interface {
	Store(ctx context.Context, ev CrawlEvent) (CrawlEvent, error)
	Get(ctx context.Context, id int) (CrawlEvent, error)
	HasRecentlyCrawled(ctx context.Context, u string) (bool, error)
}

type DbCachedCrawlerRepository struct {
	db    *bun.DB
	cache *Cache
}

func NewDbCachedCrawlerRepository(db *bun.DB, cache *Cache) CrawlerRepository {
	return &DbCachedCrawlerRepository{db, cache}
}

func (r *DbCachedCrawlerRepository) Store(ctx context.Context, ev CrawlEvent) (CrawlEvent, error) {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().
			Model(&ev).
			Exec(ctx)

		if err != nil {
			return err
		}

		for i := range ev.Headings {
			ev.Headings[i].CrawlId = ev.Id
		}
		for i := range ev.Links {
			ev.Links[i].CrawlId = ev.Id
		}

		_, err = tx.NewInsert().
			Model(&ev.Headings).
			Exec(ctx)

		if err != nil {
			return err
		}

		_, err = tx.NewInsert().
			Model(&ev.Links).
			Exec(ctx)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return CrawlEvent{}, err
	}

	log.Printf("Stored CrawlEvent: %d", ev.Id)

	return ev, nil
}

func (r *DbCachedCrawlerRepository) Get(ctx context.Context, id int) (CrawlEvent, error) {
	crawlEvent := new(CrawlEvent)
	err := r.db.NewSelect().
		Model(crawlEvent).
		Relation("Headings").
		Relation("Links").
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return *crawlEvent, err
	}

	return *crawlEvent, nil
}

func (r *DbCachedCrawlerRepository) HasRecentlyCrawled(ctx context.Context, u string) (bool, error) {
	hit := r.cache.HasRecentlyCrawled(ctx, u)

	if hit {
		return hit, nil
	}

	crawlEvent := new(CrawlEvent)
	err := r.db.NewSelect().
		Model(crawlEvent).
		Where("url = ?", u).
		Scan(ctx)

	if err != nil {
		return false, err
	}

	return true, nil
}
