package repositories

import (
	"fmt"
	"time"

	"github.com/greatfocus/gf-frame/cache"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-user/models"
)

// PersonRepository struct
type PersonRepository struct {
	db    *database.Conn
	cache *cache.Cache
}

// Init method
func (repo *PersonRepository) Init(db *database.Conn, cache *cache.Cache) {
	repo.db = db
	repo.cache = cache
}

// Create method
func (repo *PersonRepository) Create(person models.Person) (models.Person, error) {
	statement := `
    INSERT INTO person (userId, countryId, firstName, middleName, lastName, mobileNumber, idNumber, createdBy, updatedBy)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    returning id
  `
	var id int64
	err := repo.db.Master.Conn.QueryRow(statement, person.UserID, person.CountryID, person.FirstName, person.MiddleName, person.LastName, person.MobileNumber, person.IDNumber, person.UserID, person.UserID).Scan(&id)
	if err != nil {
		return person, err
	}
	createdPerson := person
	createdPerson.ID = id
	return createdPerson, nil
}

// GetByUserID method
func (repo *PersonRepository) GetByUserID(userID int64) (models.Person, error) {
	// get data from cache
	var key = "PersonRepository.GetByUserID" + string(userID)
	found, cache := repo.getPersonCache(key)
	if found {
		return cache, nil
	}

	query := `
	SELECT id, userId, countryId, firstName, middleName, lastName, mobileNumber, idNumber
	FROM person
	WHERE userId = $1
    `
	row := repo.db.Slave.Conn.QueryRow(query, userID)
	var person models.Person
	err := row.Scan(&person.ID, &person.UserID, &person.CountryID, &person.FirstName, &person.MiddleName, &person.LastName, &person.MobileNumber, &person.IDNumber)
	if err != nil {
		return models.Person{}, err
	}

	// update cache
	repo.setPersonCache(key, person)
	return person, nil
}

// Update method
func (repo *PersonRepository) Update(person models.Person) error {
	query := `
    UPDATE person
	SET 
		countryId=$2,
		firstName=$3,
		middleName=$4,
		lastName=$5,
		mobileNumber=$6,
		idNumber=$7,
		updatedOn=CURRENT_TIMESTAMP
    WHERE id=$1
  	`

	res, err := repo.db.Master.Conn.Exec(query, person.ID, person.CountryID, person.FirstName,
		person.MiddleName, person.LastName, person.MobileNumber, person.IDNumber)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got Update Person for %d", person.ID)
	}

	return nil
}

// getPersonCache method get cache for person
func (repo *PersonRepository) getPersonCache(key string) (bool, models.Person) {
	var data models.Person
	if x, found := repo.cache.Get(key); found {
		data = x.(models.Person)
		return found, data
	}
	return false, data
}

// setPersonCache method set cache for person
func (repo *PersonRepository) setPersonCache(key string, person models.Person) {
	if person != (models.Person{}) {
		repo.cache.Set(key, person, 5*time.Minute)
	}
}
