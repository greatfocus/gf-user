package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/greatfocus/gf-frame/config"
	"github.com/greatfocus/gf-frame/database"
	frameRepositories "github.com/greatfocus/gf-frame/repositories"
	"github.com/greatfocus/gf-frame/responses"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/repositories"
)

// ForgotPasswordController struct
type ForgotPasswordController struct {
	notifyRepository *frameRepositories.NotifyRepository
	userRepository   *repositories.UserRepository
	otpRepository    *repositories.OtpRepository
	config           *config.Config
}

// Init method
func (c *ForgotPasswordController) Init(db *database.DB, config *config.Config) {
	c.userRepository = &repositories.UserRepository{}
	c.userRepository.Init(db)

	c.otpRepository = &repositories.OtpRepository{}
	c.otpRepository.Init(db)

	c.notifyRepository = &frameRepositories.NotifyRepository{}
	c.notifyRepository.Init(db)

	c.config = config
}

// Handler method routes to http methods supported
func (c *ForgotPasswordController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.resetPassword(w, r)
	default:
		err := errors.New("Invalid Request")
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
}

// resetPassword method
func (c *ForgotPasswordController) resetPassword(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		derr := errors.New("invalid payload request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	user.PrepareInput()
	err = user.Validate("otp")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// check for user
	userFound, err := c.userRepository.GetByEmail(user.Email)
	if err != nil {
		derr := errors.New("User does not exist")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	if userFound.ID == 0 {
		derr := errors.New("User does not exist")
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	// activate user and verify token
	userFound.Status = "USER.CREATED"
	userFound.Enabled = false
	err = c.userRepository.UpdateLoginAttempt(userFound)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	// create new OTP
	otp := models.Otp{}
	otp.PrepareInput()
	otp.UserID = userFound.ID
	createToken, err := c.otpRepository.Create(otp)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	// create alert
	userFound.Token = createToken.Token
	if err := sendOTP(c.notifyRepository, c.config, userFound); err != nil {
		err := errors.New("unexpected error occurred")
		log.Println(err)
	}

	result := models.User{}
	result.PrepareOutput(user)
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, result.ID))
	responses.Success(w, http.StatusCreated, result)
}
