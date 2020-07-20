package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-frame/jwt"
	"github.com/greatfocus/gf-frame/responses"
	"github.com/greatfocus/gf-frame/utils"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/repositories"
)

// LoginController struct
type LoginController struct {
	userRepository  *repositories.UserRepository
	rightRepository *repositories.RightRepository
}

// Init method
func (c *LoginController) Init(db *database.DB) {
	c.userRepository = &repositories.UserRepository{}
	c.userRepository.Init(db)

	c.rightRepository = &repositories.RightRepository{}
	c.rightRepository.Init(db)
}

// Handler method routes to http methods supported
func (c *LoginController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.login(w, r)
	default:
		err := errors.New("Invalid Request")
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
}

// login method
func (c *LoginController) login(w http.ResponseWriter, r *http.Request) {
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
	err = user.Validate("login")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// check for duplicates
	userFound, err := c.userRepository.GetPasswordByEmail(user.Email)
	if err != nil {
		derr := errors.New("User does not exist or inactive")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	// check for login user status
	if userFound.Status != "USER.VERIFIED" {
		derr := errors.New("User not verified")
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	if userFound.FailedAttempts > 4 {
		userFound.Status = "USER.LOCKED"
		userFound.Enabled = false
		err = c.userRepository.UpdateLoginAttempt(userFound)

		derr := errors.New("User account is locked")
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	// check for login attempts
	if userFound.ID == 0 {
		derr := errors.New("User does not exist")
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	// verify password
	userFound.LastAttempt = time.Now()
	userFound.LastChange = time.Now()
	var valid = utils.ComparePasswords(userFound.Password, []byte(user.Password))
	if !valid {
		derr := errors.New("username of password is invalid")
		log.Printf("Error: %v\n", derr)
		userFound.FailedAttempts = (userFound.FailedAttempts + 1)
		c.userRepository.UpdateLoginAttempt(userFound)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	// update attempts
	userFound.FailedAttempts = 0
	err = c.userRepository.UpdateLoginAttempt(userFound)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	// append user rights
	right, err := c.rightRepository.GetRight(userFound.ID)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	user.Right = right

	// generate token
	token, err := jwt.CreateToken(userFound.ID, right.Role)
	user.JWT = token
	result := models.User{}
	result.PrepareOutput(user)
	responses.Success(w, http.StatusCreated, result)
}
