package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/greatfocus/gf-sframe/database"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/repositories"
	cache "github.com/patrickmn/go-cache"
)

// OtpService struct
type OtpService struct {
	userRepository *repositories.UserRepository
	otpRepository  *repositories.OtpRepository
}

// Init method
func (o *OtpService) Init(database database.Database, cache *cache.Cache) {
	o.userRepository = &repositories.UserRepository{}
	o.userRepository.Init(database, cache)

	o.otpRepository = &repositories.OtpRepository{}
	o.otpRepository.Init(database, cache)
}

// ValidateToken method
func (o *OtpService) ValidateToken(ctx context.Context, enKey string, otp models.Otp) (models.Otp, error) {
	// get token from DB
	insertedOtp, err := o.otpRepository.GetByToken(ctx, enKey, otp.Token)
	if err != nil {
		derr := errors.New("token invalid")
		log.Printf("Error: %v\n", err)
		return otp, derr
	}

	if insertedOtp.ExpiredDate.Before(time.Now()) {
		derr := errors.New("token expired")
		log.Printf("Error: %v\n", err)
		return otp, derr
	}

	if insertedOtp.Channel != otp.Channel || !insertedOtp.Active || insertedOtp.Verified {
		derr := errors.New("invalid token")
		log.Printf("Error: %v\n", err)
		return otp, derr
	}

	insertedOtp.Verified = true
	insertedOtp.Active = false
	err = o.otpRepository.Update(ctx, insertedOtp)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		return otp, derr
	}

	return insertedOtp, nil
}

// CreateToken provides new otp
func (o *OtpService) CreateToken(ctx context.Context, otp models.Otp, enKey string) (models.Otp, error) {
	// create new OTP
	otp.PrepareInput()
	return o.otpRepository.Create(ctx, enKey, otp)
}
