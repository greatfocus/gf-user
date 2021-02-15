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

// ForgotPasswordController struct
type ForgotPasswordController struct {
	userService *services.UserService
}

// Init method
func (f *ForgotPasswordController) Init(userService *services.UserService) {
	f.userService = userService
}

// Handler method routes to http methods supported
func (f *ForgotPasswordController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		f.resetPassword(w, r)
	default:
		err := errors.New("invalid request")
		responses.Error(w, http.StatusNotFound, err)
		return
	}
}

// resetPassword method
func (f *ForgotPasswordController) resetPassword(w http.ResponseWriter, r *http.Request) {
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
	err = user.PrepareInput()
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = user.Validate("otp")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// reset password
	result, err := f.userService.ResetPassword(user)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.Success(w, http.StatusOK, result)
}
