package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	server "github.com/greatfocus/gf-sframe/server"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/services"
)

// Login struct
type Login struct {
	LoginHandler func(http.ResponseWriter, *http.Request)
	userService  *services.UserService
	meta         *server.Meta
}

// Init method
func (l *Login) Init(meta *server.Meta, userService *services.UserService) {
	l.userService = userService
	l.meta = meta
}

func (l Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		l.login(w, r)
		return
	}

	// catch all
	// if no method is satisfied return an error
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Add("Allow", "GET, POST, PUT, DELETE")
}

// login methodhandler
func (l *Login) login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		l.meta.Error(w, r, derr)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		l.meta.Error(w, r, derr)
		return
	}
	err = user.Validate("login")
	if err != nil {
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		l.meta.Error(w, r, err)
		return
	}

	// login user
	result, err := l.userService.Login(ctx, user)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		l.meta.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	l.meta.Success(w, r, result)
}
