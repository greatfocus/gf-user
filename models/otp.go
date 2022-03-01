package models

import (
	"errors"
	"time"

	"github.com/greatfocus/gf-sframe/util"
)

// Otp struct
type Otp struct {
	ID          int64     `json:"-"`
	Token       int64     `json:"token,omitempty"`
	Channel     string    `json:"channel,omitempty"`
	ExpiredDate time.Time `json:"-"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Verified    bool      `json:"verified,omitempty"`
	Active      bool      `json:"active,omitempty"`
}

// PrepareInput initiliazes the User request object
func (u *Otp) PrepareInput() {
	// All users have expiry date of 15 minutes if they don't login
	var expire = time.Now().Add(time.Minute * 15)

	// random token
	var token = util.RandNumber(5)
	u.Token = token
	u.ExpiredDate = expire
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.Verified = false
}

// PrepareOutput initiliazes the User request object
func (u *Otp) PrepareOutput(otp Otp) {
	u.Verified = otp.Verified
	u.Token = otp.Token
}

// Validate check if request is valid
func (u *Otp) ExistingTokenValidate() error {
	if u.Token == 0 {
		return errors.New("required token")
	}
	if u.Channel == "" {
		return errors.New("requited channel")
	}
	return nil
}

// Validate check if request is valid
func (u *Otp) NewTokenValidate() error {
	if u.Channel == "" {
		return errors.New("requited channel")
	}
	return nil
}
