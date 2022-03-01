package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	server "github.com/greatfocus/gf-sframe/server"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/services"
)

// Otp struct
type Otp struct {
	OtpHandler func(http.ResponseWriter, *http.Request)
	otpService *services.OtpService
	server     *server.Server
}

// Init method
func (o *Otp) Init(server *server.Server, otpService *services.OtpService) {
	o.otpService = otpService
	o.server = server
}

// ValidateRequest checks if request is valid
func (o Otp) ValidateRequest(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	data, err := o.server.Request(w, r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (o Otp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		data, err := o.ValidateRequest(w, r)
		if err == nil {
			o.createToken(w, r, data)
		}
		return
	}

	if r.Method == http.MethodPut {
		data, err := o.ValidateRequest(w, r)
		if err == nil {
			o.validateToken(w, r, data)
		}
		return
	}

	// catch all
	// if no method is satisfied return an error
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Add("Allow", "POST, PUT")
}

// validateToken method adds a new token
func (o *Otp) createToken(w http.ResponseWriter, r *http.Request, data interface{}) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(o.server.Timeout)*time.Second)
	defer cancel()

	// validate if json object
	body, err := json.Marshal(data)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		o.server.Error(w, r, derr)
		return
	}
	otp := models.Otp{}
	err = json.Unmarshal(body, &otp)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		o.server.Error(w, r, derr)
		return
	}
	// validate payload rules
	err = otp.NewTokenValidate()
	if err != nil {
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		o.server.Error(w, r, err)
		return
	}

	// validate token
	createdOtp, err := o.otpService.CreateToken(ctx, otp, o.server.JWT.Secret())
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		o.server.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	o.server.Success(w, r, createdOtp)
}

// validateToken method check if token is valid
func (o *Otp) validateToken(w http.ResponseWriter, r *http.Request, data interface{}) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(o.server.Timeout)*time.Second)
	defer cancel()

	// validate if json object
	body, err := json.Marshal(data)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		o.server.Error(w, r, derr)
		return
	}
	otp := models.Otp{}
	err = json.Unmarshal(body, &otp)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		o.server.Error(w, r, derr)
		return
	}
	// validate payload rules
	err = otp.ExistingTokenValidate()
	if err != nil {
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		o.server.Error(w, r, err)
		return
	}

	// validate token
	createdOtp, err := o.otpService.ValidateToken(ctx, o.server.JWT.Secret(), otp)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		o.server.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	o.server.Success(w, r, createdOtp)
}
