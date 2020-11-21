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

// ClientAuthController struct
type ClientAuthController struct {
	clientService *services.ClientService
}

// Init method
func (l *ClientAuthController) Init(clientService *services.ClientService) {
	l.clientService = clientService
}

// Handler method routes to http methods supported
func (l *ClientAuthController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		l.authenticate(w, r)
	default:
		err := errors.New("Invalid Request")
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
}

// authenticate method
func (l *ClientAuthController) authenticate(w http.ResponseWriter, r *http.Request) {
	// '98590c398a254d2898838e1b17381575', 'ADRtjWLkttBbMQLpMADF'
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	client := models.Client{}
	err = json.Unmarshal(body, &client)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	err = client.Validate("auth")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// authenticate client
	result, err := l.clientService.Authenticate(client)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.Success(w, http.StatusCreated, result)
}
