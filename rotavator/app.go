package rotavator

import (
	"net/http"

	"github.com/gorilla/mux"
)

type App struct{}

func (a App) Run() error {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)

	server := http.Server{Addr: "0.0.0.0:5001", Handler: r}
	return server.ListenAndServe()
}
