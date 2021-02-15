package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/greatfocus/gf-frame/responses"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/services"
)

// UserController struct
type UserController struct {
	userService *services.UserService
}

// Init method
func (c *UserController) Init(userService *services.UserService) {
	c.userService = userService
}

// Handler method routes to http methods supported
func (c *UserController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.createUser(w, r)
	case http.MethodGet:
		c.getUsers(w, r)
	default:
		err := errors.New("invalid request")
		responses.Error(w, http.StatusNotFound, err)
		return
	}
}

// createUser method
func (c *UserController) createUser(w http.ResponseWriter, r *http.Request) {
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

	user, err = c.userService.CreateUser(user)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.Success(w, http.StatusOK, user)
	return
}

// getUsers method
func (c *UserController) getUsers(w http.ResponseWriter, r *http.Request) {
	var err error
	lastID := int64(1)
	lastIDStr := r.FormValue("lastId")
	id := int64(0)
	idStr := r.FormValue("id")

	if len(idStr) != 0 {
		id, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			derr := errors.New("invalid parameter")
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusBadRequest, derr)
			return
		}

		user := models.User{}
		user, err = c.userService.GetUser(id)
		if err != nil {
			responses.Error(w, http.StatusUnprocessableEntity, err)
			return
		}
		responses.Success(w, http.StatusOK, user)
		return
	}
	if len(lastIDStr) != 0 {
		lastID, err = strconv.ParseInt(lastIDStr, 10, 64)
		if err != nil {
			derr := errors.New("invalid parameter")
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusBadRequest, derr)
			return
		}

		users := []models.User{}
		users, err = c.userService.GetUsers(lastID)
		if err != nil {
			responses.Error(w, http.StatusUnprocessableEntity, err)
			return
		}
		responses.Success(w, http.StatusOK, users)
		return
	}

	derr := errors.New("invalid parameter")
	responses.Error(w, http.StatusBadRequest, derr)
	return
}
