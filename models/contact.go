package models

import (
	"errors"
	"strings"
	"time"
)

// Contact struct
type Contact struct {
	ID        int64     `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Message   string    `json:"message,omitempty"`
	Status    string    `json:"status,omitempty"`
	CreatedOn time.Time `json:"-"`
	UpatedOn  time.Time `json:"-"`
}

// Validate check if request is valid
func (c *Contact) Validate(action string) error {
	switch strings.ToLower(action) {
	case "contact":
		if c.Name == "" {
			return errors.New("required eassword")
		}
		if c.Email == "" {
			return errors.New("required email")
		}
		if c.Message == "" {
			return errors.New("required email")
		}
		return nil
	default:
		return errors.New("invalid payload request")
	}
}

// PrepareOutput initiliazes the contact request object
func (c *Contact) PrepareOutput(contact Contact) {
	c.ID = contact.ID
}
