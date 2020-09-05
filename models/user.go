package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/greatfocus/gf-frame/utils"
)

// User struct
type User struct {
	ID             int64     `json:"id,omitempty"`
	Type           string    `json:"type,omitempty"`
	Email          string    `json:"email,omitempty"`
	Password       string    `json:"password,omitempty"`
	JWT            string    `json:"jwt,omitempty"`
	Token          int64     `json:"-"`
	FailedAttempts int64     `json:"-"`
	LastAttempt    time.Time `json:"-"`
	ExpiredDate    time.Time `json:"-"`
	CreatedOn      time.Time `json:"-"`
	UpdatedOn      time.Time `json:"-"`
	Status         string    `json:"-"`
	Enabled        bool      `json:"-"`
	Right          Right     `json:"right,omitempty"`
}

// PrepareInput initiliazes the User request object
func (u *User) PrepareInput() {
	// All users have expiry date of 3 months if they don't login
	var expire = time.Now()
	expire.AddDate(0, 3, 0)
	var pass = utils.HashAndSalt([]byte(u.Password))
	fmt.Println(pass)

	u.ID = 0
	u.Password = pass

	u.FailedAttempts = 0
	u.LastAttempt = time.Now()
	u.ExpiredDate = expire
	u.CreatedOn = time.Now()
	u.UpdatedOn = time.Now()
	u.Enabled = false
	u.Status = "USER.CREATED"
}

// PrepareOutput initiliazes the User request object
func (u *User) PrepareOutput(user User) {
	u.ID = user.ID
	u.Type = user.Type
	u.Email = user.Email
	u.JWT = user.JWT
	u.Right = user.Right
}

// Validate check if request is valid
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "otp":
		if u.Email == "" {
			return errors.New("Required Email")
		}
		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		return nil
	case "edit":
		if u.ID == 0 {
			return errors.New("Required id")
		}
		if u.Type == "" {
			return errors.New("Required Type")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		return nil
	case "register":
		if u.Type == "" {
			return errors.New("Required Type")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		return nil
	default:
		return errors.New("invalid payload request")
	}
}
