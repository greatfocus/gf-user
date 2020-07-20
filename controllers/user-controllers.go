package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/greatfocus/gf-frame/config"
	"github.com/greatfocus/gf-frame/database"
	frameRepositories "github.com/greatfocus/gf-frame/repositories"
	"github.com/greatfocus/gf-frame/responses"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/repositories"
)

// UserController struct
type UserController struct {
	userRepository   *repositories.UserRepository
	otpRepository    *repositories.OtpRepository
	rightRepository  *repositories.RightRepository
	notifyRepository *frameRepositories.NotifyRepository
	config           *config.Config
}

// Init method
func (c *UserController) Init(db *database.DB, config *config.Config) {
	c.userRepository = &repositories.UserRepository{}
	c.userRepository.Init(db)

	c.otpRepository = &repositories.OtpRepository{}
	c.otpRepository.Init(db)

	c.rightRepository = &repositories.RightRepository{}
	c.rightRepository.Init(db)

	c.notifyRepository = &frameRepositories.NotifyRepository{}
	c.notifyRepository.Init(db)

	c.config = config
}

// Handler method routes to http methods supported
func (c *UserController) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.createUser(w, r)
	case http.MethodPut:
		c.updateUser(w, r)
	case http.MethodGet:
		c.getUsers(w, r)
	default:
		err := errors.New("Invalid Request")
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
}

// createUser method
func (c *UserController) createUser(w http.ResponseWriter, r *http.Request) {
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
	err = user.Validate("")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// check for duplicates
	users := []models.User{}
	users, err = c.userRepository.GetByEmailOrMobileNumber(user.Email, user.MobileNumber)
	if err != nil {
		derr := errors.New("User already exist")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}
	if len(users) > 0 {
		derr := errors.New("User already exist")
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	createdUser, err := c.userRepository.CreateUser(user)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	// create default role
	right := models.Right{}
	right.UserID = createdUser.ID
	_, err = c.rightRepository.CreateDefault(right)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	// create new OTP
	otp := models.Otp{}
	otp.PrepareInput()
	otp.UserID = createdUser.ID
	createToken, err := c.otpRepository.Create(otp)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	// create alert
	createdUser.Token = createToken.Token
	if err := sendOTP(c.notifyRepository, c.config, createdUser); err != nil {
		err := errors.New("unexpected error occurred")
		log.Println(err)
	}

	result := models.User{}
	result.PrepareOutput(createdUser)
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, result.ID))
	responses.Success(w, http.StatusCreated, result)
}

// getUsers method
func (c *UserController) getUsers(w http.ResponseWriter, r *http.Request) {
	var err error
	page := int64(1)
	pageStr := r.FormValue("page")
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

		user := models.User{}
		user, err := c.userRepository.GetUser(id)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, err)
			return
		}
		responses.Success(w, http.StatusOK, user)
		return
	}
	if len(pageStr) != 0 {
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			derr := errors.New("Invalid parameter")
			log.Printf("Error: %v\n", err)
			responses.Error(w, http.StatusBadRequest, derr)
			return
		}

		users := []models.User{}
		users, err = c.userRepository.GetUsers(page)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, err)
			return
		}
		responses.Success(w, http.StatusOK, users)
		return
	}

	derr := errors.New("Invalid parameter")
	responses.Error(w, http.StatusBadRequest, derr)
	return
}

// createUser method
func (c *UserController) updateUser(w http.ResponseWriter, r *http.Request) {
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
	err = user.Validate("edit")
	if err != nil {
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = c.userRepository.UpdateUser(user)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		responses.Error(w, http.StatusUnprocessableEntity, derr)
		return
	}

	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, user.ID))
	responses.Success(w, http.StatusCreated, user)
}

// sendOTP create alerts
func sendOTP(repo *frameRepositories.NotifyRepository, c *config.Config, user models.User) error {
	output := make([]string, 2)
	output[0] = user.FirstName
	output[1] = strconv.Itoa(int(user.Token))
	err := repo.SendNotification(c, output, user.Email, user.ID, "email_otp")
	if err != nil {
		return err
	}
	return nil
}
