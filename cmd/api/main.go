package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"vincent-h-lee/web-crawler/internal/api"
	"vincent-h-lee/web-crawler/internal/crawler"
	"vincent-h-lee/web-crawler/internal/queue"
	"vincent-h-lee/web-crawler/internal/util"

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
	log.Print("Running")

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
	publisher := queue.NewRabbitMqPublisher(&q, ch)

	sqldb := sql.OpenDB(getPgDriverConnector())
	db := bun.NewDB(sqldb, pgdialect.New())

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("CACHE_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	cache := crawler.NewCache(redisClient)

	app := api.NewApp(":8080", db, publisher, cache)
	app.Start()
}
