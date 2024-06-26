package crawler

import (
	"context"
	"database/sql"
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
	cache Cache
}

func NewDbCachedCrawlerRepository(db *bun.DB, cache Cache) CrawlerRepository {
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

	r.cache.SetUrl(ctx, ev.Url)

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
	hit, err := r.cache.HasRecentlyCrawled(ctx, u)

	if err != nil {
		return false, err
	}

	if hit {
		return true, nil
	}

	crawlEvent := new(CrawlEvent)
	err = r.db.NewSelect().
		Model(crawlEvent).
		Where("url = ?", u).
		Where("timestamp > now() - interval '1 day'").
		Scan(ctx)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	// result found so store in cache
	r.cache.SetUrl(ctx, u)

	return true, nil
}
