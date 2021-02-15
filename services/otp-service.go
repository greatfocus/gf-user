package services

import (
	"errors"
	"log"
	"time"

	"github.com/greatfocus/gf-frame/server"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/repositories"
)

// OtpService struct
type OtpService struct {
	userRepository *repositories.UserRepository
	otpRepository  *repositories.OtpRepository
}

// Init method
func (o *OtpService) Init(s *server.Server) {
	o.userRepository = &repositories.UserRepository{}
	o.userRepository.Init(s.DB, s.Cache)

	o.otpRepository = &repositories.OtpRepository{}
	o.otpRepository.Init(s.DB, s.Cache)
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
		derr := errors.New("token invalid")
		log.Printf("Error: %v\n", err)
		return otp, derr
	}

	if dbOtp.ID == 0 {
		derr := errors.New("token invalid")
		log.Printf("Error: %v\n", err)
		return otp, derr
	}

	dbOtp.ExpiredDate.Add(time.Minute * 30)
	if dbOtp.ExpiredDate.Before(time.Now()) {
		derr := errors.New("token expired")
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
