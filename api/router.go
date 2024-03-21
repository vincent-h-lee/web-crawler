package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type approuter struct {
	Router *chi.Mux
}

func NewRouter() *approuter {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world."))
	})

	return &approuter{
		Router: router,
	}
}
