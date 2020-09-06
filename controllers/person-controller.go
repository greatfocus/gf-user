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

// PersonController struct
type PersonController struct {
	personService *services.PersonService
}

// Init method
func (p *PersonController) Init(personService *services.PersonService) {
	p.personService = personService
}

// Handler method routes to http methods supported
func (p *PersonController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p.Create(w, r)
	case http.MethodPut:
		p.Update(w, r)
	case http.MethodGet:
		p.Get(w, r)
	default:
		err := errors.New("Invalid Request")
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
}

// validateToken method
func (p *PersonController) Create(w http.ResponseWriter, r *http.Request) {
	// Get body from request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	// validate if json object
	person := models.Person{}
	err = json.Unmarshal(body, &person)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	// validate payload rules
	err = person.Validate("create")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// create token
	createdPerson, err := p.personService.Create(person)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.Success(w, http.StatusCreated, createdPerson)
}

// Update method
func (p *PersonController) Update(w http.ResponseWriter, r *http.Request) {
	// Get body from request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	// validate if json object
	person := models.Person{}
	err = json.Unmarshal(body, &person)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	// validate payload rules
	err = person.Validate("create")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// create token
	createdPerson, err := p.personService.Update(person)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.Success(w, http.StatusCreated, createdPerson)
}

// Get method
func (p *PersonController) Get(w http.ResponseWriter, r *http.Request) {
	var err error
	userID := int64(0)
	userIDStr := r.FormValue("userId")

	if len(userIDStr) != 0 {
		userID, err = strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusBadRequest, derr)
			return
		}

		person := models.Person{}
		person, err := p.personService.Get(userID)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, err)
			return
		}
		responses.Success(w, http.StatusOK, person)
	} else {
		derr := errors.New("Invalid parameter")
		responses.Error(w, http.StatusBadRequest, derr)
	}
}
