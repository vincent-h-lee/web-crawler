package main

import (
	"log"
	"vincent-h-lee/web-crawler/api"
)

func main() {
	log.Print("Running")
	app := api.NewApp(":8080")
	app.Start()
}
