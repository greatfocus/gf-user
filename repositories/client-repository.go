package repositories

import (
	"database/sql"
	"fmt"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-user/models"
)

// ClientRepository struct
type ClientRepository struct {
	db *database.DB
}

// Init method
func (repo *ClientRepository) Init(db *database.DB) {
	repo.db = db
}

// Create method
func (repo *ClientRepository) Create(client models.Client) (models.Client, error) {
	statement := `
    insert into client (email, clientId, secret, expiredDate)
    values ($1, $2, $3, $4)
    returning id
  `
	var id int64
	err := repo.db.Conn.QueryRow(statement, client.Email, client.ClientID, client.Secret, client.ExpiredDate).Scan(&id)
	if err != nil {
		return client, err
	}
	created := client
	created.ID = id
	return created, nil
}

// GetByEmail method
func (repo *ClientRepository) GetByEmail(email string) (models.Client, error) {
	var client models.Client
	query := `
	select id, email, failedAttempts, lastAttempt, expiredDate, createdOn, updatedOn, enabled
	from client
	where email = $1 and deleted=false
    `
	row := repo.db.Conn.QueryRow(query, email)
	err := row.Scan(&client.ID, &client.Email, &client.FailedAttempts,
		&client.LastAttempt, &client.ExpiredDate, &client.CreatedOn, &client.UpdatedOn, &client.Enabled)
	if err != nil {
		return models.Client{}, err
	}

	return client, nil
}

// Login method
func (repo *ClientRepository) Login(client models.Client) (models.Client, error) {
	query := `
	select id, email, failedAttempts, lastAttempt, expiredDate, createdOn, updatedOn, enabled
	from client
	where email = $1 and clientId = $2 and secret = $3
    `
	row := repo.db.Conn.QueryRow(query, client.Email, client.ClientID, client.Secret)
	err := row.Scan(&client.ID, &client.Email, &client.FailedAttempts,
		&client.LastAttempt, &client.ExpiredDate, &client.CreatedOn, &client.UpdatedOn, &client.Enabled)
	if err != nil {
		return models.Client{}, err
	}

	return client, nil
}

// Delete method
func (repo *ClientRepository) Delete(id int64) error {
	query := `
    delete from client
    where id=$1
  	`
	res, err := repo.db.Conn.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got updated client for %d", id)
	}

	return nil
}

// UpdateLoginAttempt method
func (repo *ClientRepository) UpdateLoginAttempt(client models.Client) error {
	query := `
    update client
	set 
		lastAttempt=$2, 
		failedAttempts=$3,
		enabled=$4		
    where id=$1
  	`

	res, err := repo.db.Conn.Exec(query, client.ID, client.LastAttempt, client.FailedAttempts, client.Enabled)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got updated client for %d", client.ID)
	}

	return nil
}

// GetClients method
func (repo *ClientRepository) GetClients(page int64) ([]models.Client, error) {
	query := `
	select id, email, failedAttempts, lastAttempt, expiredDate, createdOn, updatedOn, enabled
	from client
	order BY createdOn limit 50 OFFSET $1-1
    `
	rows, err := repo.db.Conn.Query(query, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getClientsFromRows(rows)
}

// GetClient method
func (repo *ClientRepository) GetClient(id int64) (models.Client, error) {
	var client models.Client
	query := `
	select id, email, failedAttempts, lastAttempt, expiredDate, createdOn, updatedOn, enabled
	from client
	where id=$1 and enabled=true
	`
	row := repo.db.Conn.QueryRow(query, id)
	err := row.Scan(&client.ID, &client.Email, &client.FailedAttempts, &client.LastAttempt,
		&client.ExpiredDate, &client.CreatedOn, &client.UpdatedOn, &client.Enabled)
	if err != nil {
		return models.Client{}, err
	}

	return client, nil
}

// prepare clients row
func getClientsFromRows(rows *sql.Rows) ([]models.Client, error) {
	clients := []models.Client{}
	for rows.Next() {
		var client models.Client
		err := rows.Scan(&client.ID, &client.Email, &client.FailedAttempts, &client.LastAttempt,
			&client.ExpiredDate, &client.CreatedOn, &client.UpdatedOn, &client.Enabled)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}
