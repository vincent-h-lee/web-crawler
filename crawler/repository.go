package crawler

import (
	"context"

	"github.com/uptrace/bun"
)

type CrawlerRepository struct {
	db *bun.DB
}

func NewCrawlerRepository(db *bun.DB) *CrawlerRepository {
	return &CrawlerRepository{db}
}

func (r *CrawlerRepository) Store(ctx context.Context, ev CrawlEvent) (CrawlEvent, error) {
	var newEv CrawlEvent
	var id int64

	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewInsert().
			Model(&ev).
			Exec(ctx)

		if err != nil {
			return err
		}

		id, _ = res.LastInsertId()

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

	err = r.db.NewSelect().
		Model(newEv).
		Relation("Headings").
		Relation("Links").
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return CrawlEvent{}, err
	}

	return ev, nil
}

func (r *CrawlerRepository) Get(ctx context.Context, id int) (CrawlEvent, error) {
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
