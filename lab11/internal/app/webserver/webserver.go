package webserver

import (
	"net/http"
)

func Start(config *Config) error {
	s := newServer()

	return http.ListenAndServe(config.BindAddr, s)
}
