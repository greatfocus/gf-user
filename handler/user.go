package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	server "github.com/greatfocus/gf-sframe/server"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/services"
)

// User struct
type User struct {
	UserHandler func(http.ResponseWriter, *http.Request)
	userService *services.UserService
	meta        *server.Meta
}

// Init method
func (c *User) Init(meta *server.Meta, userService *services.UserService) {
	c.userService = userService
	c.meta = meta
}

func (c User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.createUser(w, r)
		return
	} else if r.Method == http.MethodGet {
		c.getUsers(w, r)
		return
	}

	// catch all
	// if no method is satisfied return an error
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Add("Allow", "GET, POST, PUT, DELETE")
}

// createUser method
func (c *User) createUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		c.meta.Error(w, r, derr)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		c.meta.Error(w, r, derr)
		return
	}

	user, err = c.userService.CreateUser(ctx, user)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		c.meta.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	c.meta.Success(w, r, user)
}

// getUsers method
func (c *User) getUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

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
			w.WriteHeader(http.StatusBadRequest)
			c.meta.Error(w, r, derr)
			return
		}

		user := models.User{}
		user, err = c.userService.GetUser(ctx, id)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			c.meta.Error(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		c.meta.Success(w, r, user)
		return
	}
	if len(lastIDStr) != 0 {
		lastID, err = strconv.ParseInt(lastIDStr, 10, 64)
		if err != nil {
			derr := errors.New("invalid parameter")
			log.Printf("Error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			c.meta.Error(w, r, derr)
			return
		}

		users, err := c.userService.GetUsers(ctx, lastID)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			c.meta.Error(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		c.meta.Success(w, r, users)
		return
	}

	derr := errors.New("invalid parameter")
	c.meta.Error(w, r, derr)
	w.WriteHeader(http.StatusBadRequest)
}
