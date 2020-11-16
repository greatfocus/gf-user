package services

import (
	"errors"
	"log"
	"time"

	"github.com/greatfocus/gf-frame/server"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/repositories"
)

// ClientService struct
type ClientService struct {
	clientRepository *repositories.ClientRepository
}

// Init method
func (u *ClientService) Init(s *server.Server) {
	u.clientRepository = &repositories.ClientRepository{}
	u.clientRepository.Init(s.DB)
}

// Create method
func (u *ClientService) Create(client models.Client) (models.Client, error) {
	client.PrepareInput()
	err := client.Validate("create")
	if err != nil {
		derr := errors.New("Invalid request")
		log.Printf("Error: %v\n", err)
		return client, derr
	}

	// check for duplicates
	exist := models.Client{}
	exist, err = u.clientRepository.GetByEmail(client.Email)
	if (models.Client{}) != exist {
		derr := errors.New("Client already exist")
		log.Printf("Error: %v\n", err)
		return client, derr
	}

	// Create client
	created, err := u.clientRepository.Create(client)
	if err != nil {
		derr := errors.New("Client registration failed")
		log.Printf("Error: %v\n", err)
		return client, derr
	}

	result := models.Client{}
	result.PrepareOutput(created)
	return result, nil
}

// GetByID method
func (u *ClientService) GetByID(id int64) (models.Client, error) {
	client, err := u.clientRepository.GetClient(id)
	if err != nil {
		derr := errors.New("User does not exist")
		log.Printf("Error: %v\n", err)
		return client, derr
	}
	return client, nil
}

// GetClients method
func (u *ClientService) GetClients(page int64) ([]models.Client, error) {
	clients, err := u.clientRepository.GetClients(page)
	if err != nil {
		derr := errors.New("Failed to fetch User")
		log.Printf("Error: %v\n", err)
		return clients, derr
	}
	return clients, nil
}

// Authenticate method
func (u *ClientService) Authenticate(client models.Client) (models.Client, error) {
	// check for duplicates
	found, err := u.clientRepository.Login(client)
	if err != nil {
		derr := errors.New("Client does not exist or inactive")
		log.Printf("Error: %v\n", err)
		return client, derr
	}

	if found.FailedAttempts > 4 {
		found.Enabled = false
		err = u.clientRepository.UpdateLoginAttempt(found)

		derr := errors.New("Client account is locked")
		log.Printf("Error: %v\n", err)
		return client, derr
	}
	// check for login attempts
	if found.ID == 0 {
		derr := errors.New("Client does not exist")
		log.Printf("Error: %v\n", derr)
		return client, derr
	}

	// verify password
	found.LastAttempt = time.Now()
	if found.ClientID != client.ClientID && found.Secret != client.Secret {
		derr := errors.New("Client ID or Secret is invalid")
		log.Printf("Error: %v\n", derr)
		found.FailedAttempts = (found.FailedAttempts + 1)
		u.clientRepository.UpdateLoginAttempt(found)
		return client, derr
	}

	// update attempts
	found.FailedAttempts = 0
	err = u.clientRepository.UpdateLoginAttempt(found)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		return client, derr
	}

	result := models.Client{}
	result.PrepareOutput(client)
	return result, nil
}

// Delete method
func (u *ClientService) Delete(id int64) error {
	// check for user
	found, err := u.clientRepository.GetClient(id)
	if err != nil {
		derr := errors.New("User does not exist")
		log.Printf("Error: %v\n", err)
		return derr
	}
	if found.ID == 0 {
		derr := errors.New("User does not exist")
		log.Printf("Error: %v\n", derr)
		return derr
	}

	err = u.clientRepository.Delete(found.ID)
	if err != nil {
		derr := errors.New("unexpected error occurred. kindly initiate forget password request")
		log.Printf("Error: %v\n", err)
		return derr
	}

	return nil
}
