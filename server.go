package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vasart/go-rest-api/model"
)

type Server struct {
	router *mux.Router
}

func NewServer(u model.UserRepository) *Server {
	s := Server{router: mux.NewRouter()}
	NewUserRouter(u, s.newSubrouter("/user"))
	return &s
}

func (s *Server) Start() {
	log.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", s.router); err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
}

func (s *Server) newSubrouter(path string) *mux.Router {
	return s.router.PathPrefix(path).Subrouter()
}
