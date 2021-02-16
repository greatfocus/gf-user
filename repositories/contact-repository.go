package repositories

import (
	"fmt"

	"github.com/greatfocus/gf-frame/cache"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-user/models"
)

// ContactRepository struct
type ContactRepository struct {
	db    *database.Conn
	cache *cache.Cache
}

// Init method
func (repo *ContactRepository) Init(db *database.Conn, cache *cache.Cache) {
	repo.db = db
	repo.cache = cache
}

// ReachToUs method
func (repo *ContactRepository) ReachToUs(contact models.Contact) (models.Contact, error) {
	statement := `
    INSERT INTO contactus (name, email, message, status, createdOn)
    VALUES ($1, $2, $3, 'new', $4)
    returning id
  `
	var id int64
	err := repo.db.Master.Conn.QueryRow(statement, contact.Name, contact.Email, contact.Message, contact.CreatedOn).Scan(&id)
	if err != nil {
		return contact, err
	}
	created := contact
	created.ID = id
	return created, nil
}

// GetMessages method
func (repo *ContactRepository) GetMessages(status string) ([]models.Contact, error) {
	query := `
	SELECT name, email, message, status, createdOn, updatedOn
	FROM contactus 
	WHERE status = $1
	ORDER BY id DESC LIMIT 20
    `
	rows, err := repo.db.Slave.Conn.Query(query, status)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

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
    UPDATE contactus
	SET 
		status=$2,
		updatedOn=CURRENT_TIMESTAMP
    WHERE id=$1
  	`

	res, err := repo.db.Master.Conn.Exec(query, contact.ID, "new")
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
