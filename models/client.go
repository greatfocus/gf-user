package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/greatfocus/gf-frame/utils"
)

// Client struct
type Client struct {
	ID             int64     `json:"id,omitempty"`
	Email          string    `json:"email,omitempty"`
	ClientID       string    `json:"clientId,omitempty"`
	Secret         string    `json:"secret,omitempty"`
	ClientIDTmp    string    `json:"-"`
	SecretTmp      string    `json:"-"`
	FailedAttempts int64     `json:"-"`
	LastAttempt    time.Time `json:"-"`
	ExpiredDate    time.Time `json:"-"`
	CreatedOn      time.Time `json:"-"`
	UpdatedOn      time.Time `json:"-"`
	Enabled        bool      `json:"-"`
	Authenticated  bool      `json:"authenticated,omitempty"`
}

// PrepareInput initiliazes the Client request object
func (u *Client) PrepareInput() error {
	// All users have expiry date of 3 months if they don't login
	var expire = time.Now()
	expire.AddDate(0, 3, 0)

	// client id
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	clientID, err := utils.HashAndSalt([]byte(uuid))
	if err != nil {
		return errors.New("Failed to create hash")
	}
	u.ClientID = clientID
	u.ClientIDTmp = uuid

	// secret
	var rand = utils.RandomString(20)
	secret, err := utils.HashAndSalt([]byte(rand))
	if err != nil {
		return errors.New("Failed to create hash")
	}
	u.Secret = secret
	u.SecretTmp = rand

	u.FailedAttempts = 0
	u.LastAttempt = time.Now()
	u.ExpiredDate = expire
	u.CreatedOn = time.Now()
	u.UpdatedOn = time.Now()
	u.Enabled = false

	return nil
}

// PrepareOutput initiliazes the Client request object
func (u *Client) PrepareOutput(client Client) {
	u.Email = client.Email
	u.Authenticated = client.Authenticated
}

// Validate check if request is valid
func (u *Client) Validate(action string) error {
	switch strings.ToLower(action) {
	case "create":
		if u.Email == "" {
			return errors.New("Required Email")
		}
		return nil
	case "delete":
		if u.ID == 0 {
			return errors.New("Required ID")
		}
		if u.ClientID == "" {
			return errors.New("Required Client ID")
		}
		if u.Secret == "" {
			return errors.New("Required Secret")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		return nil
	case "auth":
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if u.ClientID == "" {
			return errors.New("Required Client ID")
		}
		if u.Secret == "" {
			return errors.New("Required Secret")
		}
		return nil
	default:
		return errors.New("invalid payload request")
	}
}
