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

// OtpController struct
type OtpController struct {
	otpService *services.OtpService
}

// Init method
func (o *OtpController) Init(otpService *services.OtpService) {
	o.otpService = otpService
}

// Handler method routes to http methods supported
func (o *OtpController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		o.validateToken(w, r)
	default:
		err := errors.New("invalid request")
		responses.Error(w, http.StatusNotFound, err)
		return
	}
}

// validateToken method
func (o *OtpController) validateToken(w http.ResponseWriter, r *http.Request) {
	// Get body from request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusBadRequest, derr)
		return
	}
	// validate if json object
	otp := models.Otp{}
	err = json.Unmarshal(body, &otp)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusBadRequest, derr)
		return
	}
	// validate payload rules
	err = otp.Validate()
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// validate token
	createdOtp, err := o.otpService.ValidateToken(otp)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.Success(w, http.StatusOK, createdOtp)
	return
}
