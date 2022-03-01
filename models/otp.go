package models

import (
	"errors"
	"time"

	"github.com/greatfocus/gf-sframe/crypt"
)

// Otp struct
type Otp struct {
	ID          int64     `json:"id,omitempty"`
	UserID      int64     `json:"userid,omitempty"`
	Email       string    `json:"email,omitempty"`
	Token       int64     `json:"token,omitempty"`
	ExpiredDate time.Time `json:"-"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Verified    bool      `json:"-"`
}

// PrepareInput initiliazes the User request object
func (u *Otp) PrepareInput() {
	// All users have expiry date of 30 minutes if they don't login
	var expire = time.Now()
	expire.Add(time.Minute * 30)

	// random token
	var token = crypt.RandomNumber(5)

	u.ID = 0
	u.UserID = 0
	u.Token = token
	u.ExpiredDate = expire
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.Verified = false
}

// PrepareOutput initiliazes the User request object
func (u *Otp) PrepareOutput(otp Otp) {
	u.Verified = otp.Verified
}

// Validate check if request is valid
func (u *Otp) Validate() error {

	if u.Token == 0 {
		return errors.New("required token")
	}
	if u.Email == "" {
		return errors.New("required email")
	}
	return nil
}
