package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"vincent-h-lee/web-crawler/internal/crawler"
	"vincent-h-lee/web-crawler/internal/queue"
	"vincent-h-lee/web-crawler/internal/util"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func getPgDriverConnector() *pgdriver.Connector {
	var (
		host     = os.Getenv("POSTGRES_HOST")
		port     = os.Getenv("POSTGRES_PORT")
		user     = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbname   = os.Getenv("POSTGRES_DB")
	)

	return pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", host, port)),
		pgdriver.WithDatabase(dbname),
		pgdriver.WithUser(user),
		pgdriver.WithPassword(password),
		pgdriver.WithInsecure(true),
	)
}

func main() {
	conn, err := amqp.Dial(os.Getenv("MQ_CONNECTION"))
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"urls", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	util.FailOnError(err, "Failed to register a consumer")

	sqldb := sql.OpenDB(getPgDriverConnector())
	db := bun.NewDB(sqldb, pgdialect.New())

	l := launcher.MustNewManaged(os.Getenv("BROWSER_URL"))
	l.Headless(false).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")
	browser := rod.New().Client(l.MustClient()).MustConnect()
	pool := rod.NewPagePool(1) // TODO concurrency

	// Create a page if needed
	create := func() *rod.Page {
		// We use MustIncognito to isolate pages with each other
		return browser.MustIncognito().MustPage()
	}
	rodPool := crawler.NewRodPool(&pool, create)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("CACHE_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	cache := crawler.NewRedisCache(redisClient)

	repo := crawler.NewDbCachedCrawlerRepository(db, cache)
	publisher := queue.NewRabbitMqPublisher(&q, ch)

	job := queue.NewConsumer(repo, publisher, rodPool)

	var forever chan struct{}

	go func() {
		for d := range msgs {
			ctx := context.Background()
			log.Printf("Received a message: %s", d.Body)
			err := job(ctx, string(d.Body))
			if err != nil {
				log.Printf("Processing message failed with error: %s", err)
			} else {
				log.Printf("Processed message")
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
