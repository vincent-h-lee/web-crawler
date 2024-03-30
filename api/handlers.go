package api

import (
	"net/http"
	"strconv"
	"time"
	"vincent-h-lee/web-crawler/crawler"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	repo *crawler.CrawlerRepository
}

func (h *Handlers) GetCrawlEventRoute(w http.ResponseWriter, r *http.Request) {
	crawlIdParam := chi.URLParam(r, "crawlId")
	crawlId, err := strconv.Atoi(crawlIdParam)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid crawl event ID")
	}

	ev, err := h.repo.Get(r.Context(), crawlId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, ev)
}

func (h *Handlers) PostCrawlEventRoute(w http.ResponseWriter, r *http.Request) {
	ev, err := h.repo.Store(r.Context(), crawler.NewCrawlEvent("https://millandcommonsss.com", "", []crawler.Heading{{Text: "some title", Tag: "h1"}}, []crawler.Link{{Text: "Link", Url: "https://google.com"}}, 400, time.Now()))
	/* consumer := queue.NewConsumer(repo)
	ev, err := consumer.Consume("https://millandcommons.com") */

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, ev)
}
