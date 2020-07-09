package models

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/greatfocus/gf-frame/utils"
)

// User struct
type User struct {
	ID             int64     `json:"id,omitempty"`
	Type           string    `json:"type,omitempty"`
	FirstName      string    `json:"firstname,omitempty"`
	MiddleName     string    `json:"middlename,omitempty"`
	LastName       string    `json:"lastname,omitempty"`
	MobileNumber   string    `json:"mobilenumber,omitempty"`
	Email          string    `json:"email,omitempty"`
	Password       string    `json:"password,omitempty"`
	JWT            string    `json:"jwt,omitempty"`
	Token          int64     `json:"-"`
	FailedAttempts int64     `json:"failedattempts,omitempty"`
	LastAttempt    time.Time `json:"-"`
	LastChange     time.Time `json:"-"`
	ExpiredDate    time.Time `json:"expireddate,omitempty"`
	CreatedOn      time.Time `json:"-"`
	CreatedBy      time.Time `json:"-"`
	UpdatedOn      time.Time `json:"-"`
	Status         string    `json:"status,omitempty"`
	Enabled        bool      `json:"enabled,omitempty"`
	Rights         []Right   `json:"rights,omitempty"`
}

// PrepareInput initiliazes the User request object
func (u *User) PrepareInput() {
	// All users have expiry date of 3 months if they don't login
	var expire = time.Now()
	expire.AddDate(0, 3, 0)
	var pass = utils.HashAndSalt([]byte(u.Password))
	fmt.Println(pass)

	u.ID = 0
	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
	u.MiddleName = html.EscapeString(strings.TrimSpace(u.MiddleName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
	u.MobileNumber = html.EscapeString(strings.TrimSpace(u.MobileNumber))
	u.Password = pass

	u.FailedAttempts = 0
	u.LastAttempt = time.Now()
	u.LastChange = time.Now()
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
	u.FirstName = user.FirstName
	u.MiddleName = user.MiddleName
	u.LastName = user.LastName
	u.MobileNumber = user.MobileNumber
	u.Email = user.Email
	u.JWT = user.JWT
	u.Rights = user.Rights
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
		if u.FirstName == "" {
			return errors.New("Required First Name")
		}
		if u.MiddleName == "" {
			return errors.New("Required Middle Name")
		}
		if u.LastName == "" {
			return errors.New("Required Last Name")
		}
		if u.MobileNumber == "" {
			return errors.New("Required Mobile Number")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		return nil
	default:
		if u.Type == "" {
			return errors.New("Required Type")
		}
		if u.FirstName == "" {
			return errors.New("Required First Name")
		}
		if u.MiddleName == "" {
			return errors.New("Required Middle Name")
		}
		if u.LastName == "" {
			return errors.New("Required Last Name")
		}
		if u.MobileNumber == "" {
			return errors.New("Required Mobile Number")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		return nil
	}
}
