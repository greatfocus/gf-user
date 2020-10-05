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

// ContactController struct
type ContactController struct {
	userService *services.UserService
}

// Init method
func (l *ContactController) Init(userService *services.UserService) {
	l.userService = userService
}

// Handler method routes to http methods supported
func (l *ContactController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		l.reachToUs(w, r)
	default:
		err := errors.New("Invalid Request")
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
}

// reachToUs method
func (l *ContactController) reachToUs(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	user := models.Contact{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	err = user.Validate("contact")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// reachToUs user
	result, err := l.userService.reachToUs(user)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.Success(w, http.StatusCreated, result)
}
