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

func (r *CrawlerRepository) Store(ev CrawlEvent) (CrawlEvent, error) {
	// TODO
	return CrawlEvent{}, nil
}

func (r *CrawlerRepository) Get(ctx context.Context, id int) (CrawlEvent, error) {
	crawlEvent := new(CrawlEvent)
	err := r.db.NewSelect().Model(crawlEvent).Relation("Headings").Relation("Links").Where("id = ?", 1).Scan(ctx)
	if err != nil {
		return *crawlEvent, err
	}
	return *crawlEvent, nil
}
