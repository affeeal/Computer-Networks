package webserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newServer() *server {
	s := &server{
		router: mux.NewRouter(),
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/dashboard", s.handleDashboard()).Methods("GET")
}

func (s *server) handleDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/dashboard.html")
	}
}
