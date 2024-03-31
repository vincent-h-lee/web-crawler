package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"vincent-h-lee/web-crawler/internal/crawler"
	"vincent-h-lee/web-crawler/internal/queue"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	publisher *queue.Publisher
	repo      *crawler.CrawlerRepository
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

type PostCrawlEventBody struct {
	Url string `json:"url"`
}

func (h *Handlers) PostCrawlEventRoute(w http.ResponseWriter, r *http.Request) {
	var body PostCrawlEventBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if _, err := url.ParseRequestURI(body.Url); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid URL")
		return
	}

	err := h.publisher.Publish(r.Context(), body.Url)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
