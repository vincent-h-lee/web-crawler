package api

import (
	"log"
	"net/http"

	"github.com/uptrace/bun"
)

type app struct {
	srv *http.Server
}

func (a *app) Start() {
	log.Printf("Starting server at port %s", a.srv.Addr)
	log.Fatal(a.srv.ListenAndServe())
}

func NewApp(addr string, db *bun.DB) *app {
	routes := NewRouter(db)

	srv := &http.Server{
		Addr:    addr,
		Handler: routes.Router,
	}

	return &app{srv}
}
