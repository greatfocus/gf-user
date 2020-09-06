package repositories

import (
	"fmt"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-user/models"
)

// PersonRepository struct
type PersonRepository struct {
	db *database.DB
}

// Init method
func (repo *PersonRepository) Init(db *database.DB) {
	repo.db = db
}

// Create method
func (repo *PersonRepository) Create(person models.Person) (models.Person, error) {
	statement := `
    insert into person (userId, countryId, firstName, middleName, lastName, mobileNumber, idNumber, createdBy, updatedBy)
    values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    returning id
  `
	var id int64
	err := repo.db.Conn.QueryRow(statement, person.UserID, person.CountryID, person.FirstName, person.MiddleName, person.LastName, person.MobileNumber, person.IdNumber, person.UserID, person.UserID).Scan(&id)
	if err != nil {
		return person, err
	}
	createdPerson := person
	createdPerson.ID = id
	return createdPerson, nil
}

// GetByUserId method
func (repo *PersonRepository) GetByUserId(userID int64) (models.Person, error) {
	query := `
	select id, userId, countryId, firstName, middleName, lastName, mobileNumber, idNumber
	from person
	where userId = $1
    `
	row := repo.db.Conn.QueryRow(query, userID)
	var person models.Person
	err := row.Scan(&person.ID, &person.UserID, &person.CountryID, &person.FirstName, &person.MiddleName, &person.LastName, &person.MobileNumber, &person.IdNumber)
	if err != nil {
		return models.Person{}, err
	}

	return person, nil
}

// Update method
func (repo *PersonRepository) Update(person models.Person) error {
	query := `
    update person
	set 
		countryId=$2,
		firstName=$3,
		middleName=$4,
		lastName=$5,
		mobileNumber=$6,
		idNumber=$7,
		updatedOn=CURRENT_TIMESTAMP
    where id=$1
  	`

	res, err := repo.db.Conn.Exec(query, person.ID, person.CountryID, person.FirstName,
		person.MiddleName, person.LastName, person.MobileNumber, person.IdNumber)
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
