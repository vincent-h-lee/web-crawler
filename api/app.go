package api

import (
	"log"
	"net/http"
)

type app struct {
	srv *http.Server
}

func (a *app) Start() {
	log.Printf("Starting server at port %s", a.srv.Addr)
	log.Fatal(a.srv.ListenAndServe())
}

func NewApp(addr string) *app {
	/* dbService, err := service.NewDatabaseService("./enforcements.db")
	if err != nil {
		log.Fatal(err)
	} */

	routes := NewRouter()

	srv := &http.Server{
		Addr:    addr,
		Handler: routes.Router,
	}

	return &app{srv}
}
