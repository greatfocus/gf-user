package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/greatfocus/gf-frame/responses"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/services"
)

// LoginController struct
type LoginController struct {
	userService *services.UserService
}

// Init method
func (l *LoginController) Init(userService *services.UserService) {
	l.userService = userService
}

// Handler method routes to http methods supported
func (l *LoginController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		l.login(w, r)
	default:
		err := errors.New("invalid request")
		responses.Error(w, http.StatusNotFound, err)
		return
	}
}

// login method
func (l *LoginController) login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusBadRequest, derr)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusBadRequest, derr)
		return
	}
	err = user.Validate("login")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// login user
	result, err := l.userService.Login(user)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.Success(w, http.StatusOK, result)
}
