package apiserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/affeeal/lab9/internal/app/model"
	"github.com/affeeal/lab9/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

const (
	sessionName = "session"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type server struct {
	router       *mux.Router
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/signup", s.handleSignupGet()).Methods("GET")
	s.router.HandleFunc("/signup", s.handleSignupPost()).Methods("POST")

	s.router.HandleFunc("/login", s.handleLoginGet()).Methods("GET")
	s.router.HandleFunc("/login", s.handleLoginPost()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)

	internal := s.router.PathPrefix("/internal").Subrouter()

	private.HandleFunc("/sync", s.handlePrivateSync()).Methods("GET")
	internal.HandleFunc("/sync", s.handleInternalSync()).Methods("POST")

	private.HandleFunc("/async", s.handlePrivateAsync()).Methods("GET")
	internal.HandleFunc("/async", s.handleInternalAsync()).Methods("GET")
}

func (s *server) handleSignupGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/signup.html")
	}
}

func (s *server) handleSignupPost() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleLoginGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/login.html")
	}
}

func (s *server) handleLoginPost() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err = s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		_, err = s.store.User().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) handlePrivateSync() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/sync.html")
	}
}

func (s *server) handleInternalSync() http.HandlerFunc {
	type message struct {
		Out string `json:"out"`
		Err string `json:"err"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		msg := r.PostForm.Get("message")
		buffout, bufferr := s.executeCommand(msg)
		s.respond(w, r, http.StatusOK, message{Out: buffout.String(), Err: bufferr.String()})
	}
}

func (s *server) handlePrivateAsync() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/async.html")
	}
}

func (s *server) handleInternalAsync() http.HandlerFunc {
	type message struct {
		Out string `json:"out"`
		Err string `json:"err"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		defer c.Close()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			buffout, bufferr := s.executeCommand(string(msg))
			res, err := json.Marshal(message{Out: buffout.String(), Err: bufferr.String()})
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			if err = c.WriteMessage(websocket.TextMessage, res); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
		}
	}
}

func (s *server) executeCommand(msg string) (bytes.Buffer, bytes.Buffer) {
	split := strings.Split(msg, " ")
	cmd := exec.Command(split[0], split[1:]...)

	var buffout, bufferr bytes.Buffer
	cmd.Stdout = &buffout
	cmd.Stderr = &bufferr
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}

	log.Println(buffout.String())

	return buffout, bufferr
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
