package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
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
	meta       *server.Meta
}

// Init method
func (o *Otp) Init(meta *server.Meta, otpService *services.OtpService) {
	o.otpService = otpService
	o.meta = meta
}

func (o Otp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		o.validateToken(w, r)
		return
	}

	// catch all
	// if no method is satisfied return an error
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Add("Allow", "GET, POST, PUT, DELETE")
}

// validateToken method
func (o *Otp) validateToken(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Get body from request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		o.meta.Error(w, r, derr)
		return
	}
	// validate if json object
	otp := models.Otp{}
	err = json.Unmarshal(body, &otp)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		o.meta.Error(w, r, derr)
		return
	}
	// validate payload rules
	err = otp.Validate()
	if err != nil {
		log.Printf("Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		o.meta.Error(w, r, err)
		return
	}

	// validate token
	createdOtp, err := o.otpService.ValidateToken(ctx, otp)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		o.meta.Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	o.meta.Success(w, r, createdOtp)
}
