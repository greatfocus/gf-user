package repositories

import (
	"fmt"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-user/models"
)

// ContactRepository struct
type ContactRepository struct {
	db *database.DB
}

// Init method
func (repo *ContactRepository) Init(db *database.DB) {
	repo.db = db
}

// ReachToUs method
func (repo *ContactRepository) ReachToUs(contact models.Contact) (models.Contact, error) {
	statement := `
    insert into contactus (name, email, message, status, createdOn)
    values ($1, $2, $3, 'new', $4)
    returning id
  `
	var id int64
	err := repo.db.Conn.QueryRow(statement, contact.Name, contact.Email, contact.Message, contact.CreatedOn).Scan(&id)
	if err != nil {
		return contact, err
	}
	createdContact := contact
	createdContact.ID = id
	return createdContact, nil
}

// GetMessages method
func (repo *ContactRepository) GetMessages(status string) ([]models.Contact, error) {
	query := `
	select name, email, message, status, createdOn, updatedOn
	from contactus 
	where status = $1
	order BY createdOn limit 10
    `
	rows, err := repo.db.Conn.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []models.Contact{}
	for rows.Next() {
		var message models.Contact
		err := rows.Scan(&message.ID, &message.Name, &message.Email, &message.Message, &message.Status)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// Update method
func (repo *ContactRepository) Update(contact models.Contact) error {
	query := `
    update contactus
	set 
		status=$2,
		updatedOn=CURRENT_TIMESTAMP
    where id=$1
  	`

	res, err := repo.db.Conn.Exec(query, contact.ID, "new")
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got Update Contact for %d", contact.ID)
	}

	return nil
}
