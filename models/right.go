package models

import (
	"errors"
	"time"
)

// Right struct
type Right struct {
	ID        int64     `json:"-"`
	UserID    int64     `json:"-"`
	RoleID    int64     `json:"roleId,omitempty"`
	Role      string    `json:"role,omitempty"`
	ActionID  string    `json:"actionId,omitempty"`
	Action    string    `json:"action,omitempty"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Deleted   bool      `json:"-"`
	Enabled   bool      `json:"enabled,omitempty"`
}

// PrepareInput initiliazes the User request object
func (u *Right) PrepareInput() {
	u.ID = 0
	u.UserID = 0
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.Deleted = false
	u.Enabled = true
}

// PrepareOutput initiliazes the User request object
func (u *Right) PrepareOutput(right Right) {
	u.UserID = right.UserID
	u.ID = right.ID
}

// Validate check if request is valid
func (u *Right) Validate() error {

	if u.UserID == 0 {
		return errors.New("required user id")
	}
	return nil
}
