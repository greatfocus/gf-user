package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/greatfocus/gf-sframe/crypt"
	"github.com/greatfocus/gf-sframe/database"
	"github.com/greatfocus/gf-sframe/server"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/repositories"
	cache "github.com/patrickmn/go-cache"
)

// UserService struct
type UserService struct {
	userRepository  *repositories.UserRepository
	otpRepository   *repositories.OtpRepository
	rightRepository *repositories.RightRepository
	jwt             server.JWT
}

// Init method
func (u *UserService) Init(database database.Database, cache *cache.Cache, jwt server.JWT) {
	u.userRepository = &repositories.UserRepository{}
	u.userRepository.Init(database, cache)

	u.otpRepository = &repositories.OtpRepository{}
	u.otpRepository.Init(database, cache)

	u.rightRepository = &repositories.RightRepository{}
	u.rightRepository.Init(database, cache)

	u.jwt = jwt
}

// CreateUser method
func (u *UserService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	err := user.PrepareInput(u.jwt.Secret())
	if err != nil {
		return user, err
	}
	err = user.Validate("register")
	if err != nil {
		derr := errors.New("invalid request")
		log.Printf("Error: %v\n", err)
		return user, derr
	}

	// isValid := validate.Email(user.Email)
	// if !isValid {
	// 	derr := errors.New("invalid email address")
	// 	log.Printf("Error: %v\n", derr)
	// 	return user, derr
	// }

	// check for duplicates
	usersExist, err := u.userRepository.GetByEmail(ctx, user.Email)
	if (models.User{}) != usersExist {
		derr := errors.New("user already exist")
		log.Printf("Error: %v\n", err)
		return user, derr
	}

	// Create user
	createdUser, err := u.userRepository.CreateUser(ctx, user)
	if err != nil {
		derr := errors.New("user registration failed")
		log.Printf("Error: %v\n", err)
		return user, derr
	}

	// create default role
	right := models.Right{}
	right.UserID = createdUser.ID
	_, err = u.rightRepository.CreateDefault(ctx, right)
	if err != nil {
		derr := errors.New("user registration failed")
		log.Printf("Error: %v\n", err)
		_ = u.userRepository.Delete(ctx, createdUser.ID)
		return user, derr
	}

	result := models.User{}
	result.PrepareOutput(createdUser)
	return result, nil
}

// GetUser method
func (u *UserService) GetUser(ctx context.Context, id int64) (models.User, error) {
	user, err := u.userRepository.GetUser(ctx, id)
	if err != nil {
		derr := errors.New("user does not exist")
		log.Printf("Error: %v\n", err)
		return user, derr
	}
	return user, nil
}

// GetUsers method
func (u *UserService) GetUsers(ctx context.Context, lastID int64) ([]models.User, error) {
	user, err := u.userRepository.GetUsers(ctx, lastID)
	if err != nil {
		derr := errors.New("failed to fetch user")
		log.Printf("Error: %v\n", err)
		return user, derr
	}
	return user, nil
}

// Login method
func (u *UserService) Login(ctx context.Context, user models.User, origin string) (models.User, error) {
	// check for duplicates
	userFound, err := u.userRepository.GetPasswordByEmail(ctx, user.Email)
	if err != nil {
		derr := errors.New("user does not exist or inactive")
		log.Printf("Error: %v\n", err)
		return user, derr
	}
	// check for login user status
	if userFound.Status != "USER.VERIFIED" && userFound.Status != "USER.APPROVED" {
		derr := errors.New("user not verified")
		log.Printf("Error: %v\n", derr)
		return user, derr
	}
	if userFound.FailedAttempts > 4 {
		userFound.Status = "USER.LOCKED"
		userFound.Enabled = false
		err = u.userRepository.UpdateLoginAttempt(ctx, userFound)

		derr := errors.New("user account is locked")
		log.Printf("Error: %v\n", err)
		return user, derr
	}
	// check for login attempts
	if userFound.ID == 0 {
		derr := errors.New("user does not exist")
		log.Printf("Error: %v\n", derr)
		return user, derr
	}

	// verify password
	userFound.LastAttempt = time.Now()
	if crypt.Decrypt(userFound.Password, u.jwt.Secret()) != user.Password {
		derr := errors.New("username of password is invalid")
		log.Printf("Error: %v\n", derr)
		userFound.FailedAttempts = (userFound.FailedAttempts + 1)
		_ = u.userRepository.UpdateLoginAttempt(ctx, userFound)
		return user, derr
	}

	// update attempts
	userFound.FailedAttempts = 0
	userFound.SuccessLogins = userFound.SuccessLogins + 1
	err = u.userRepository.UpdateLoginAttempt(ctx, userFound)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		return user, derr
	}

	// append user rights
	right, err := u.rightRepository.GetRight(ctx, userFound.ID)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		return user, derr
	}
	user.Right = right

	// send first time login message
	// if userFound.SuccessLogins == 1 {
	// 	_ = sendFirstTimeLogin(u.notifyRepository, u.config, userFound)
	// }

	// generate token

	tokenInfo := server.TokenInfo{
		ActorID: userFound.ID,
		Origin:  origin,
	}
	token, _ := u.jwt.CreateToken(tokenInfo)
	user.JWT = token
	result := models.User{}
	result.PrepareOutput(user)
	return result, nil
}

// ResetPassword method
func (u *UserService) ResetPassword(ctx context.Context, user models.User) (models.User, error) {
	// check for user
	userFound, err := u.userRepository.GetByEmail(ctx, user.Email)
	if err != nil {
		derr := errors.New("user does not exist")
		log.Printf("Error: %v\n", err)
		return user, derr
	}
	if userFound.ID == 0 {
		derr := errors.New("user does not exist")
		log.Printf("Error: %v\n", derr)
		return user, derr
	}

	// activate user and verify token
	userFound.Status = "USER.CREATED"
	userFound.Enabled = false
	err = u.userRepository.UpdateLoginAttempt(ctx, userFound)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		return user, derr
	}

	// create new OTP
	otp := models.Otp{}
	otp.PrepareInput()
	// createToken, err := u.otpRepository.Create(otp, "email")
	// if err != nil {
	// 	derr := errors.New("unexpected error occurred")
	// 	log.Printf("Error: %v\n", err)
	// 	return user, derr
	// }

	// create alert
	// userFound.Token = createToken.Token
	// if err := sendOTP(u.notifyRepository, u.config, userFound); err != nil {
	// 	err := errors.New("unexpected error occurred")
	// 	log.Println(err)
	// }

	//_ = sendResetPassword(u.notifyRepository, u.config, userFound)

	result := models.User{}
	result.PrepareOutput(user)
	return result, nil
}

// // sendResetPassword create alerts
// func sendResetPassword(repo *frameRepositories.NotifyRepository, c *config.Config, user models.User) error {
// 	output := make([]string, 0)
// 	err := repo.AddNotification(c, output, user.Email, user.ID, "password_reset")
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // sendOTP create alerts
// func sendOTP(repo *frameRepositories.NotifyRepository, c *config.Config, user models.User) error {
// 	output := make([]string, 1)
// 	output[0] = strconv.Itoa(int(user.Token))
// 	err := repo.AddNotification(c, output, user.Email, user.ID, "email_otp")
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // sendFirstTimeLogin create alerts
// func sendFirstTimeLogin(repo *frameRepositories.NotifyRepository, c *config.Config, user models.User) error {
// 	output := make([]string, 0)
// 	err := repo.AddNotification(c, output, user.Email, user.ID, "first_login")
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
