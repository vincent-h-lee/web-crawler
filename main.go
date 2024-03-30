package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"vincent-h-lee/web-crawler/api"

	"github.com/uptrace/bun"

	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func getConnector() *pgdriver.Connector {
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

	sqldb := sql.OpenDB(getConnector())
	db := bun.NewDB(sqldb, pgdialect.New())

	app := api.NewApp(":8080", db)
	app.Start()
}
