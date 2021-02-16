package services

import (
	"errors"
	"log"

	"github.com/greatfocus/gf-frame/validate"

	"github.com/greatfocus/gf-frame/utils"

	"github.com/greatfocus/gf-frame/config"
	frameRepositories "github.com/greatfocus/gf-frame/repositories"
	"github.com/greatfocus/gf-frame/server"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/repositories"
)

// ClientService struct
type ClientService struct {
	clientRepository *repositories.ClientRepository
	notifyRepository *frameRepositories.NotifyRepository
	config           *config.Config
}

// Init method
func (u *ClientService) Init(s *server.Server) {
	u.clientRepository = &repositories.ClientRepository{}
	u.clientRepository.Init(s.DB, s.Cache)

	u.notifyRepository = &frameRepositories.NotifyRepository{}
	u.notifyRepository.Init(s.DB)

	u.config = s.Config
}

// Create method
func (u *ClientService) Create(client models.Client) (models.Client, error) {
	err := client.Validate("create")
	if err != nil {
		derr := errors.New("invalid request")
		log.Printf("Error: %v\n", err)
		return client, derr
	}

	isValid := validate.Email(client.Email)
	if !isValid {
		derr := errors.New("invalid email address")
		log.Printf("Error: %v\n", derr)
		return client, derr
	}

	// check for duplicates
	exist := models.Client{}
	exist, err = u.clientRepository.GetByEmail(client.Email)
	if (models.Client{}) != exist {
		derr := errors.New("client already exist")
		log.Printf("Error: %v\n", err)
		return client, derr
	}

	// Create client
	err = client.PrepareInput()
	if err != nil {
		log.Printf("Error: %v\n", err)
		return client, err
	}

	created, err := u.clientRepository.Create(client)
	if err != nil {
		derr := errors.New("client registration failed")
		log.Printf("Error: %v\n", err)
		return client, derr
	}

	// create alert
	if err := sendClientCredentials(u.notifyRepository, u.config, client); err != nil {
		derr := errors.New("client registration failed")
		log.Printf("Error: %v\n", err)
		u.Delete(created.ID)
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
		derr := errors.New("client does not exist")
		log.Printf("Error: %v\n", err)
		return client, derr
	}
	return client, nil
}

// GetClients method
func (u *ClientService) GetClients(lastID int64) ([]models.Client, error) {
	clients, err := u.clientRepository.GetClients(lastID)
	if err != nil {
		derr := errors.New("failed to fetch clients")
		log.Printf("Error: %v\n", err)
		return clients, derr
	}
	return clients, nil
}

// Authenticate method
func (u *ClientService) Authenticate(client models.Client) (models.Client, error) {
	// check for duplicates
	found, err := u.clientRepository.GetByEmail(client.Email)
	if err != nil {
		derr := errors.New("client does not exist or inactive")
		log.Printf("Error: %v\n", err)
		return client, derr
	}

	// check client id
	valid, err := utils.ComparePasswords(found.ClientID, []byte(client.ClientID))
	if !valid || err != nil {
		derr := errors.New("client or Secret is invalid")
		log.Printf("Error: %v\n", derr)
		found.FailedAttempts = (found.FailedAttempts + 1)
		_ = u.clientRepository.UpdateLoginAttempt(found)
		return client, derr
	}

	// check secret id
	valid, err = utils.ComparePasswords(found.Secret, []byte(client.Secret))
	if !valid || err != nil {
		derr := errors.New("client or Secret is invalid")
		log.Printf("Error: %v\n", derr)
		found.FailedAttempts = (found.FailedAttempts + 1)
		_ = u.clientRepository.UpdateLoginAttempt(found)
		return client, derr
	}

	if found.FailedAttempts > 4 {
		found.Enabled = false
		err = u.clientRepository.UpdateLoginAttempt(found)
		derr := errors.New("client account is locked")
		log.Printf("Error: %v\n", err)
		return client, derr
	}
	// check for login attempts
	if found.ID == 0 {
		derr := errors.New("client does not exist")
		log.Printf("Error: %v\n", derr)
		return client, derr
	}

	// update attempts
	found.FailedAttempts = 0
	found.Authenticated = true
	err = u.clientRepository.UpdateLoginAttempt(found)
	if err != nil {
		derr := errors.New("unexpected error occurred")
		log.Printf("Error: %v\n", err)
		return client, derr
	}

	result := models.Client{}
	result.PrepareOutput(found)
	return result, nil
}

// Delete method
func (u *ClientService) Delete(id int64) error {
	// check for user
	found, err := u.clientRepository.GetClient(id)
	if err != nil {
		derr := errors.New("client does not exist")
		log.Printf("Error: %v\n", err)
		return derr
	}
	if found.ID == 0 {
		derr := errors.New("client does not exist")
		log.Printf("Error: %v\n", derr)
		return derr
	}

	err = u.clientRepository.Delete(found.ID)
	if err != nil {
		derr := errors.New("failed to delete client")
		log.Printf("Error: %v\n", err)
		return derr
	}

	return nil
}

// sendClientCredentials create alerts
func sendClientCredentials(repo *frameRepositories.NotifyRepository, c *config.Config, client models.Client) error {
	output := make([]string, 2)
	output[0] = client.ClientIDTmp
	output[1] = client.SecretTmp
	err := repo.AddNotification(c, output, client.Email, client.ID, "client_credentials")
	if err != nil {
		return err
	}
	return nil
}
