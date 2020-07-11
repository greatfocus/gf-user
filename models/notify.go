package models

import (
	"errors"
)

// Notify struct
type Notify struct {
	ID         int64    `json:"id,omitempty"`
	UserID     int64    `json:"userId,omitempty"`
	TemplateID int64    `json:"templateId,omitempty"`
	Operation  string   `json:"operation,omitempty"`
	ChannelID  int64    `json:"channelId,omitempty"`
	Recipient  string   `json:"recipient,omitempty"`
	Param      []string `json:"param,omitempty"`
	Status     string   `json:"status,omitempty"`
	Sent       bool     `json:"sent"`
}

// PrepareNotify initiliazes the User request object
func (n *Notify) PrepareNotify() {
	n.ID = 0
	n.Status = "queue"
	n.Sent = false
}

// PrepareOutput initiliazes the User request object
func (n *Notify) PrepareOutput(notify Notify) {
	n.ID = notify.ID
	n.Sent = notify.Sent
}

// ValidateNotify check if request is valid
func (n *Notify) ValidateNotify() error {
	if n.ChannelID == 0 {
		return errors.New("Required ChannelID")
	}
	if n.Recipient == "" {
		return errors.New("Required Recipient")
	}
	if n.Recipient == "" {
		return errors.New("Required Recipient")
	}
	return nil
}
