package services

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/greatfocus/gf-frame/config"
	"github.com/greatfocus/gf-frame/jwt"
	frameRepositories "github.com/greatfocus/gf-frame/repositories"
	"github.com/greatfocus/gf-frame/server"
	"github.com/greatfocus/gf-frame/utils"
	"github.com/greatfocus/gf-frame/validate"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/repositories"
)

// UserService struct
type UserService struct {
	userRepository    *repositories.UserRepository
	otpRepository     *repositories.OtpRepository
	rightRepository   *repositories.RightRepository
	notifyRepository  *frameRepositories.NotifyRepository
	contactRepository *repositories.ContactRepository
	config            *config.Config
	jwt               *jwt.JWT
}

// Init method
func (u *UserService) Init(s *server.Server) {
	u.userRepository = &repositories.UserRepository{}
	u.userRepository.Init(s.DB, s.Cache)

	u.otpRepository = &repositories.OtpRepository{}
	u.otpRepository.Init(s.DB, s.Cache)

	u.rightRepository = &repositories.RightRepository{}
	u.rightRepository.Init(s.DB, s.Cache)

	u.notifyRepository = &frameRepositories.NotifyRepository{}
	u.notifyRepository.Init(s.DB)

	u.contactRepository = &repositories.ContactRepository{}
	u.contactRepository.Init(s.DB, s.Cache)

	u.config = s.Config
	u.jwt = s.JWT
}

// CreateUser method
func (u *UserService) CreateUser(user models.User) (models.User, error) {
	err := user.PrepareInput()
	if err != nil {
		return user, err
	}
	err = user.Validate("register")
	if err != nil {
		derr := errors.New("Invalid request")
		log.Printf("Error: %v\n", err)
		return user, derr
	}

	isValid := validate.Email(user.Email)
	if !isValid {
		derr := errors.New("Invalid Email Address")
		log.Printf("Error: %v\n", derr)
		return user, derr
	}

	// check for duplicates
	usersExist := models.User{}
	usersExist, err = u.userRepository.GetByEmail(user.Email)
	if (models.User{}) != usersExist {
		derr := errors.New("User already exist")
		log.Printf("Error: %v\n", err)
		return user, derr
	}

	// Create user
	createdUser, err := u.userRepository.CreateUser(user)
	if err != nil {
		derr := errors.New("User registration failed")
		log.Printf("Error: %v\n", err)
		return user, derr
	}

	// create default role
	right := models.Right{}
	right.UserID = createdUser.ID
	_, err = u.rightRepository.CreateDefault(right)
	if err != nil {
		derr := errors.New("User registration failed")
		log.Printf("Error: %v\n", err)
		u.userRepository.Delete(createdUser.ID)
		return user, derr
	}

	// create new OTP
	otp := models.Otp{}
	otp.PrepareInput()
	otp.UserID = createdUser.ID
	createToken, err := u.otpRepository.Create(otp, "email")
	if err != nil {
		derr := errors.New("User registration failed")
		log.Printf("Error: %v\n", err)
		u.userRepository.Delete(createdUser.ID)
		return user, derr
	}

	// create alert
	createdUser.Token = createToken.Token
	if err := sendOTP(u.notifyRepository, u.config, createdUser); err != nil {
		derr := errors.New("User registration failed")
		log.Printf("Error: %v\n", err)
		u.userRepository.Delete(createdUser.ID)
		u.otpRepository.Delete(createToken.ID)
		return user, derr
	}

	result := models.User{}
	result.PrepareOutput(createdUser)
	return result, nil
}

// GetUser method
func (u *UserService) GetUser(id int64) (models.User, error) {
	user, err := u.userRepository.GetUser(id)
	if err != nil {
		derr := errors.New("User does not exist")
		log.Printf("Error: %v\n", err)
		return user, derr
	}
	return user, nil
}

// GetUsers method
func (u *UserService) GetUsers(lastID int64) ([]models.User, error) {
	user, err := u.userRepository.GetUsers(lastID)
	if err != nil {
		derr := errors.New("Failed to fetch User")
		log.Printf("Error: %v\n", err)
		return user, derr
	}
	return user, nil
}

// Login method
func (u *UserService) Login(user models.User) (models.User, error) {
	// check for duplicates
	userFound, err := u.userRepository.GetPasswordByEmail(user.Email)
	if err != nil {
		derr := errors.New("User does not exist or inactive")
		log.Printf("Error: %v\n", err)
		return user, derr
	}
	// check for login user status
	if userFound.Status != "USER.VERIFIED" && userFound.Status != "USER.APPROVED" {
		derr := errors.New("User not verified")
		log.Printf("Error: %v\n", derr)
		return user, derr
	}
	if userFound.FailedAttempts > 4 {
		userFound.Status = "USER.LOCKED"
		userFound.Enabled = false
		err = u.userRepository.UpdateLoginAttempt(userFound)

		derr := errors.New("User account is locked")
		log.Printf("Error: %v\n", err)
		return user, derr
	}
	// check for login attempts
	if userFound.ID == 0 {
		derr := errors.New("User does not exist")
		log.Printf("Error: %v\n", derr)
		return user, derr
	}

	// verify password
	userFound.LastAttempt = time.Now()
	valid, err := utils.ComparePasswords(userFound.Password, []byte(user.Password))
	if !valid || err != nil {
		derr := errors.New("username of password is invalid")
		log.Printf("Error: %v\n", derr)
		userFound.FailedAttempts = (userFound.FailedAttempts + 1)
		u.userRepository.UpdateLoginAttempt(userFound)
		return user, derr
	}

	// update attempts
	userFound.FailedAttempts = 0
	userFound.SuccessLogins = userFound.SuccessLogins + 1
	err = u.userRepository.UpdateLoginAttempt(userFound)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		return user, derr
	}

	// append user rights
	right, err := u.rightRepository.GetRight(userFound.ID)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		return user, derr
	}
	user.Right = right

	// send first time login message
	if userFound.SuccessLogins == 1 {
		sendFirstTimeLogin(u.notifyRepository, u.config, userFound)
	}

	// generate token
	token, err := u.jwt.CreateToken(userFound.ID, right.Role)
	user.JWT = token
	result := models.User{}
	result.PrepareOutput(user)
	return result, nil
}

// ResetPassword method
func (u *UserService) ResetPassword(user models.User) (models.User, error) {
	// check for user
	userFound, err := u.userRepository.GetByEmail(user.Email)
	if err != nil {
		derr := errors.New("User does not exist")
		log.Printf("Error: %v\n", err)
		return user, derr
	}
	if userFound.ID == 0 {
		derr := errors.New("User does not exist")
		log.Printf("Error: %v\n", derr)
		return user, derr
	}

	// activate user and verify token
	userFound.Status = "USER.CREATED"
	userFound.Enabled = false
	err = u.userRepository.UpdateLoginAttempt(userFound)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		return user, derr
	}

	// create new OTP
	otp := models.Otp{}
	otp.PrepareInput()
	otp.UserID = userFound.ID
	createToken, err := u.otpRepository.Create(otp, "email")
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		return user, derr
	}

	// create alert
	userFound.Token = createToken.Token
	if err := sendOTP(u.notifyRepository, u.config, userFound); err != nil {
		err := errors.New("unexpected error occurred")
		log.Println(err)
	}

	sendResetPassword(u.notifyRepository, u.config, userFound)

	result := models.User{}
	result.PrepareOutput(user)
	return result, nil
}

// sendResetPassword create alerts
func sendResetPassword(repo *frameRepositories.NotifyRepository, c *config.Config, user models.User) error {
	output := make([]string, 0)
	err := repo.AddNotification(c, output, user.Email, user.ID, "password_reset")
	if err != nil {
		return err
	}
	return nil
}

// sendOTP create alerts
func sendOTP(repo *frameRepositories.NotifyRepository, c *config.Config, user models.User) error {
	output := make([]string, 1)
	output[0] = strconv.Itoa(int(user.Token))
	err := repo.AddNotification(c, output, user.Email, user.ID, "email_otp")
	if err != nil {
		return err
	}
	return nil
}

// sendFirstTimeLogin create alerts
func sendFirstTimeLogin(repo *frameRepositories.NotifyRepository, c *config.Config, user models.User) error {
	output := make([]string, 0)
	err := repo.AddNotification(c, output, user.Email, user.ID, "first_login")
	if err != nil {
		return err
	}
	return nil
}

// ReachToUs method
func (u *UserService) ReachToUs(contact models.Contact) (models.Contact, error) {
	err := contact.Validate("contact")
	if err != nil {
		derr := errors.New("Invalid request")
		log.Printf("Error: %v\n", err)
		return contact, derr
	}

	// Create request
	createdRequest, err := u.contactRepository.ReachToUs(contact)
	if err != nil {
		derr := errors.New("Contact request failed")
		log.Printf("Error: %v\n", err)
		return contact, derr
	}

	// create alert
	if err := sendReachToUsMessage(u.notifyRepository, u.config, createdRequest); err != nil {
		derr := errors.New("Request failed")
		log.Printf("Error: %v\n", err)
		return contact, derr
	}

	result := models.Contact{}
	result.PrepareOutput(createdRequest)
	return result, nil
}

// sendReachToUsMessage create alerts
func sendReachToUsMessage(repo *frameRepositories.NotifyRepository, c *config.Config, contact models.Contact) error {
	output := make([]string, 3)
	output[0] = contact.Name
	output[1] = contact.Email
	output[2] = contact.Message
	err := repo.AddNotification(c, output, c.Integrations.ContactUs.Email, contact.ID, "contactus_message")
	if err != nil {
		return err
	}
	return nil
}
