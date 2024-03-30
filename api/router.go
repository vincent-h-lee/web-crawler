package api

import (
	"context"
	"log"
	"net/http"
	"vincent-h-lee/web-crawler/crawler"
	"vincent-h-lee/web-crawler/queue"

	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

type approuter struct {
	Router *chi.Mux
}

func NewRouter(db *bun.DB) *approuter {
	router := chi.NewRouter()
	repo := crawler.NewCrawlerRepository(db)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		ev, err := repo.Get(ctx, 1)
		log.Print("ev", err, ev)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, ev)
	})

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		consumer := queue.NewConsumer(repo)

		ev, err := consumer.Consume("https://millandcommons.com")

		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, ev)
	})

	return &approuter{
		Router: router,
	}
}
