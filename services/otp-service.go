package services

import (
	"errors"
	"log"
	"time"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/repositories"
)

// OtpService struct
type OtpService struct {
	userRepository *repositories.UserRepository
	otpRepository  *repositories.OtpRepository
}

// Init method
func (o *OtpService) Init(db *database.DB) {
	o.userRepository = &repositories.UserRepository{}
	o.userRepository.Init(db)

	o.otpRepository = &repositories.OtpRepository{}
	o.otpRepository.Init(db)
}

// ValidateToken method
func (o *OtpService) ValidateToken(otp models.Otp) (models.Otp, error) {
	// get user via email
	user, err := o.userRepository.GetByEmail(otp.Email)
	if err != nil {
		derr := errors.New("kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		return otp, derr
	}

	// get token from DB
	dbOtp, err := o.otpRepository.GetByToken(user.ID, otp.Token)
	if err != nil {
		derr := errors.New("token Invalid")
		log.Printf("Error: %v\n", err)
		return otp, derr
	}

	if dbOtp.ID == 0 {
		derr := errors.New("token Invalid")
		log.Printf("Error: %v\n", err)
		return otp, derr
	}

	dbOtp.ExpiredDate.Add(time.Minute * 30)
	if dbOtp.ExpiredDate.Before(time.Now()) {
		derr := errors.New("Token Expired")
		log.Printf("Error: %v\n", err)
		return otp, derr
	}

	// activate user and verify token
	user.Status = "USER.VERIFIED"
	user.Enabled = true
	user.FailedAttempts = 0
	err = o.userRepository.UpdateLoginAttempt(user)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		return otp, derr
	}

	err = o.otpRepository.Update(dbOtp)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		return otp, derr
	}

	return otp, nil
}
