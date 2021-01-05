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

// ClientController struct
type ClientController struct {
	clientService *services.ClientService
}

// Init method
func (c *ClientController) Init(clientService *services.ClientService) {
	c.clientService = clientService
}

// Handler method routes to http methods supported
func (c *ClientController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.create(w, r)
	case http.MethodGet:
		c.getClients(w, r)
	case http.MethodDelete:
		c.delete(w, r)
	default:
		err := errors.New("Invalid Request")
		responses.Error(w, http.StatusNotFound, err)
		return
	}
}

// createUser method
func (c *ClientController) create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusBadRequest, derr)
		return
	}
	client := models.Client{}
	err = json.Unmarshal(body, &client)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusBadRequest, derr)
		return
	}

	client, err = c.clientService.Create(client)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.Success(w, http.StatusOK, client)
	return
}

// getClients method
func (c *ClientController) getClients(w http.ResponseWriter, r *http.Request) {
	var err error
	lastID := int64(1)
	lastIDStr := r.FormValue("lastId")
	id := int64(0)
	idStr := r.FormValue("id")

	if len(idStr) != 0 {
		id, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusBadRequest, derr)
			return
		}

		client := models.Client{}
		client, err := c.clientService.GetByID(id)
		if err != nil {
			responses.Error(w, http.StatusUnprocessableEntity, err)
			return
		}
		responses.Success(w, http.StatusOK, client)
		return
	}
	if len(lastIDStr) != 0 {
		lastID, err = strconv.ParseInt(lastIDStr, 10, 64)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusBadRequest, derr)
			return
		}

		clients := []models.Client{}
		clients, err = c.clientService.GetClients(lastID)
		if err != nil {
			responses.Error(w, http.StatusUnprocessableEntity, err)
			return
		}
		responses.Success(w, http.StatusOK, clients)
		return
	}

	derr := errors.New("Invalid parameter")
	responses.Error(w, http.StatusBadRequest, derr)
	return
}

// Delete method
func (c *ClientController) delete(w http.ResponseWriter, r *http.Request) {
	var err error
	id := int64(0)
	idStr := r.FormValue("id")

	if len(idStr) != 0 {
		id, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusBadRequest, derr)
			return
		}

		client := models.Client{}
		client.ID = id
		err := c.clientService.Delete(id)
		if err != nil {
			responses.Error(w, http.StatusUnprocessableEntity, err)
			return
		}

		responses.Success(w, http.StatusOK, client)
		return
	}
}
