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

// ForgotPassword struct
type ForgotPassword struct {
	ForgotPasswordHandler func(http.ResponseWriter, *http.Request)
	userService           *services.UserService
	server                *server.Server
}

func (f ForgotPassword) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		f.resetPassword(w, r)
		return
	}

	// catch all
	// if no method is satisfied return an error
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Add("Allow", "GET, POST, PUT, DELETE")
}

// Init method
func (f *ForgotPassword) Init(server *server.Server, userService *services.UserService) {
	f.userService = userService
	f.server = server
}

// resetPassword method
func (f *ForgotPassword) resetPassword(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(f.server.Timeout)*time.Second)
	defer cancel()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		f.server.Error(w, r, derr)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		f.server.Error(w, r, derr)
		return
	}
	err = user.PrepareInput(f.server.JWT.Secret())
	if err != nil {
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		f.server.Error(w, r, err)
		return
	}
	err = user.Validate("otp")
	if err != nil {
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		f.server.Error(w, r, err)
		return
	}

	// reset password
	result, err := f.userService.ResetPassword(ctx, user)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		f.server.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	f.server.Success(w, r, result)
}
