package queue

import (
	"context"
	"log"
	"vincent-h-lee/web-crawler/internal/crawler"

	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	repo    *crawler.CrawlerRepository
	queue   *amqp.Queue
	channel *amqp.Channel
}

func NewConsumer(repo *crawler.CrawlerRepository) *consumer {
	return &consumer{repo: repo}
}

func (c *consumer) Consume(u string) (crawler.CrawlEvent, error) {
	ctx := context.Background()

	ev, err := crawler.Crawl(u)
	if err != nil {
		log.Fatal(err)
	}

	return c.repo.Store(ctx, ev)
}
