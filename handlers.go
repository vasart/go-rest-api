package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vasart/go-rest-api/model"
)

type userRouter struct {
	userRepository model.UserRepository
}

func NewUserRouter(u model.UserRepository, router *mux.Router) *mux.Router {
	userRouter := userRouter{u}

	router.HandleFunc("/", userRouter.createUserHandler).Methods("PUT")
	router.HandleFunc("/{username}", userRouter.getUserHandler).Methods("GET")
	router.HandleFunc("/login", userRouter.loginHandler).Methods("POST")
	return router
}

func decodeUser(r *http.Request) (u model.User, err error) {
	if r.Body == nil {
		return u, errors.New("no request body")
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&u)
	return u, err
}

func decodeCredentials(r *http.Request) (c model.Credentials, err error) {
	if r.Body == nil {
		return c, errors.New("no request body")
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&c)
	return c, err
}

func (ur* userRouter) loginHandler(w http.ResponseWriter, r *http.Request) {
	credentials, err := decodeCredentials(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var user *model.User
	user, err = ur.userRepository.Login(credentials)
	if err == nil {
		Json(w, http.StatusOK, &user)
	} else {
		Error(w, http.StatusInternalServerError, "Incorrect password")
	}
}

func (ur *userRouter) createUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := decodeUser(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = ur.userRepository.CreateUser(&user)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Json(w, http.StatusOK, err)
}

func (ur *userRouter) getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars)
	username := vars["username"]

	user, err := ur.userRepository.GetByUsername(username)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}

	Json(w, http.StatusOK, user)
}
