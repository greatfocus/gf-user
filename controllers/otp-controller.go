package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-frame/responses"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/repositories"
)

// OtpController struct
type OtpController struct {
	userRepository *repositories.UserRepository
	otpRepository  *repositories.OtpRepository
}

// Init method
func (c *OtpController) Init(db *database.DB) {
	c.userRepository = &repositories.UserRepository{}
	c.userRepository.Init(db)

	c.otpRepository = &repositories.OtpRepository{}
	c.otpRepository.Init(db)
}

// Handler method routes to http methods supported
func (c *OtpController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.validateToken(w, r)
	default:
		err := errors.New("Invalid Request")
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
}

// validateToken method
func (c *OtpController) validateToken(w http.ResponseWriter, r *http.Request) {
	// Get body from request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	// validate if json object
	otp := models.Otp{}
	err = json.Unmarshal(body, &otp)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	// validate payload rules
	err = otp.Validate()
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// get user via email
	user := models.User{}
	user, err = c.userRepository.GetByEmail(otp.Email)
	if err != nil {
		derr := errors.New("kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	// get token from DB
	dbOtp, err := c.otpRepository.GetByToken(user.ID, otp.Token)
	if err != nil {
		derr := errors.New("token Invalid")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	if dbOtp.ID == 0 {
		derr := errors.New("token Invalid")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	dbOtp.ExpiredDate.Add(time.Minute * 30)
	if dbOtp.ExpiredDate.Before(time.Now()) {
		derr := errors.New("Token Expired")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	// activate user and verify token
	user.Status = "USER.VERIFIED"
	user.Enabled = true
	user.FailedAttempts = 0
	err = c.userRepository.UpdateLoginAttempt(user)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	err = c.otpRepository.Update(dbOtp)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	// Publish event

	responses.Success(w, http.StatusCreated, otp)
}
