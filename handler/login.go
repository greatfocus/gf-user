package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
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
	server       *server.Server
}

// Init method
func (l *Login) Init(server *server.Server, userService *services.UserService) {
	l.userService = userService
	l.server = server
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
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(l.server.Timeout)*time.Second)
	defer cancel()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		l.server.Error(w, r, derr)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		l.server.Error(w, r, derr)
		return
	}
	err = user.Validate("login")
	if err != nil {
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		l.server.Error(w, r, err)
		return
	}

	// login user
	result, err := l.userService.Login(ctx, user, r.Header.Get("Origin"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		l.server.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	l.server.Success(w, r, result)
}
