package services

import (
	"errors"
	"log"

	"github.com/greatfocus/gf-frame/server"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/repositories"
)

// PersonService struct
type PersonService struct {
	userRepository   *repositories.UserRepository
	personRepository *repositories.PersonRepository
}

// Init method
func (p *PersonService) Init(s *server.Server) {
	p.userRepository = &repositories.UserRepository{}
	p.userRepository.Init(s.DB)

	p.personRepository = &repositories.PersonRepository{}
	p.personRepository.Init(s.DB)
}

// Create method
func (p *PersonService) Create(person models.Person) (models.Person, error) {
	// create person details
	person, err := p.personRepository.Create(person)
	if err == nil {
		derr := errors.New("User details already exist")
		log.Printf("Error: %v\n", err)
		return person, derr
	}

	result := models.Person{}
	result.PrepareOutput(person)
	return result, nil
}

// Update method
func (p *PersonService) Update(person models.Person) (models.Person, error) {
	// create person details
	err := p.personRepository.Update(person)
	if err != nil {
		derr := errors.New("User details failed to update")
		log.Printf("Error: %v\n", err)
		return person, derr
	}

	result := models.Person{}
	result.PrepareOutput(person)
	return result, nil
}

// Get method
func (p *PersonService) Get(userID int64) (models.Person, error) {
	// create person details
	person, err := p.personRepository.GetByUserID(userID)
	if err != nil {
		derr := errors.New("User details failed to fetch")
		log.Printf("Error: %v\n", err)
		return person, derr
	}

	return person, nil
}
