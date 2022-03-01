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

// ForgotPassword struct
type ForgotPassword struct {
	ForgotPasswordHandler func(http.ResponseWriter, *http.Request)
	userService           *services.UserService
	meta                  *server.Meta
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
func (f *ForgotPassword) Init(meta *server.Meta, userService *services.UserService) {
	f.userService = userService
	f.meta = meta
}

// resetPassword method
func (f *ForgotPassword) resetPassword(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		f.meta.Error(w, r, derr)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		f.meta.Error(w, r, derr)
		return
	}
	err = user.PrepareInput()
	if err != nil {
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		f.meta.Error(w, r, err)
		return
	}
	err = user.Validate("otp")
	if err != nil {
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		f.meta.Error(w, r, err)
		return
	}

	// reset password
	result, err := f.userService.ResetPassword(ctx, user)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		f.meta.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	f.meta.Success(w, r, result)
}
