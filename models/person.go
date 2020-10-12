package models

import (
	"errors"
	"strings"
	"time"
)

// Person struct
type Person struct {
	ID           int64     `json:"id,omitempty"`
	UserID       int64     `json:"userId,omitempty"`
	CountryID    int64     `json:"countryId,omitempty"`
	FirstName    string    `json:"firstName,omitempty"`
	MiddleName   string    `json:"middleName,omitempty"`
	LastName     string    `json:"lastName,omitempty"`
	MobileNumber string    `json:"mobileNumber,omitempty"`
	IDNumber     string    `json:"idNumber,omitempty"`
	CreatedOn    time.Time `json:"-"`
	UpdatedOn    time.Time `json:"-"`
}

// PrepareInput initiliazes the Person request object
func (p *Person) PrepareInput() {
	p.ID = 0
	p.CreatedOn = time.Now()
	p.UpdatedOn = time.Now()
}

// PrepareOutput initiliazes the person request object
func (p *Person) PrepareOutput(person Person) {
	p.ID = person.ID
	p.UserID = person.UserID
}

// Validate check if request is valid
func (p *Person) Validate(action string) error {
	switch strings.ToLower(action) {
	case "create":
		if p.UserID == 0 {
			return errors.New("Required User ID")
		}
		if p.CountryID == 0 {
			return errors.New("Required Country ID")
		}
		if p.FirstName == "" {
			return errors.New("Required First Name")
		}
		if p.MiddleName == "" {
			return errors.New("Required Middle Name")
		}
		if p.LastName == "" {
			return errors.New("Required Last Name")
		}
		if p.MobileNumber == "" {
			return errors.New("Required Mobile Number")
		}
		if p.IDNumber == "" {
			return errors.New("Required ID Number")
		}
		return nil
	default:
		return errors.New("invalid payload request")
	}
}
