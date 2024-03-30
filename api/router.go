package api

import (
	"vincent-h-lee/web-crawler/crawler"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

type approuter struct {
	Router *chi.Mux
}

func NewRouter(db *bun.DB) *approuter {
	router := chi.NewRouter()
	handlers := Handlers{repo: crawler.NewCrawlerRepository(db)}

	router.Route("/crawls", func(r chi.Router) {
		r.Post("/", handlers.PostCrawlEventRoute)
		r.Get("/{crawlId}", handlers.GetCrawlEventRoute)
	})

	return &approuter{
		Router: router,
	}
}
