package api

import (
	"log"
	"net/http"
	"vincent-h-lee/web-crawler/internal/crawler"
	"vincent-h-lee/web-crawler/internal/queue"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

type app struct {
	srv *http.Server
}

func (a *app) Start() {
	log.Printf("Starting server at port %s", a.srv.Addr)
	log.Fatal(a.srv.ListenAndServe())
}

func NewApp(addr string, db *bun.DB, publisher queue.Publisher, cache crawler.Cache) *app {
	repo := crawler.NewDbCachedCrawlerRepository(db, cache)
	handlers := Handlers{repo: repo, publisher: publisher}

	router := chi.NewRouter()
	router.Route("/crawls", func(r chi.Router) {
		r.Post("/", handlers.PostCrawlEventRoute)
		r.Get("/{crawlId}", handlers.GetCrawlEventRoute)
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &app{srv}
}
