package repositories

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/greatfocus/gf-frame/cache"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-user/models"
)

// clientRepositoryCacheKeys array
var clientRepositoryCacheKeys = []string{}

// ClientRepository struct
type ClientRepository struct {
	db    *database.Conn
	cache *cache.Cache
}

// Init method
func (repo *ClientRepository) Init(db *database.Conn, cache *cache.Cache) {
	repo.db = db
	repo.cache = cache
}

// Create method
func (repo *ClientRepository) Create(client models.Client) (models.Client, error) {
	statement := `
    insert into client (email, clientId, secret, expiredDate)
    values ($1, $2, $3, $4)
    returning id
  `
	var id int64
	err := repo.db.Master.Conn.QueryRow(statement, client.Email, client.ClientID, client.Secret, client.ExpiredDate).Scan(&id)
	if err != nil {
		return client, err
	}
	created := client
	created.ID = id
	repo.deleteCache()
	return created, nil
}

// GetByEmail method
func (repo *ClientRepository) GetByEmail(email string) (models.Client, error) {
	// get data from cache
	var key = "ClientRepository.GetByEmail" + email
	found, cache := repo.getClientCache(key)
	if found {
		return cache, nil
	}

	var client models.Client
	query := `
	select id, email, clientId, secret, failedAttempts, lastAttempt, expiredDate, createdOn, updatedOn, enabled
	from client
	where email = $1 and deleted=false
    `
	row := repo.db.Slave.Conn.QueryRow(query, email)
	err := row.Scan(&client.ID, &client.Email, &client.ClientID, &client.Secret, &client.FailedAttempts,
		&client.LastAttempt, &client.ExpiredDate, &client.CreatedOn, &client.UpdatedOn, &client.Enabled)
	if err != nil {
		return models.Client{}, err
	}

	// update cache
	repo.setClientCache(key, client)
	return client, nil
}

// Login method
func (repo *ClientRepository) Login(client models.Client) (models.Client, error) {
	query := `
	select id, email, failedAttempts, lastAttempt, expiredDate, createdOn, updatedOn, enabled
	from client
	where clientId = $1 and secret = $2
    `
	row := repo.db.Slave.Conn.QueryRow(query, client.ClientID, client.Secret)
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
	res, err := repo.db.Master.Conn.Exec(query, id)
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

	repo.deleteCache()
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

	res, err := repo.db.Master.Conn.Exec(query, client.ID, client.LastAttempt, client.FailedAttempts, client.Enabled)
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

	repo.deleteCache()
	return nil
}

// GetClients method
func (repo *ClientRepository) GetClients(lastID int64) ([]models.Client, error) {
	// get data from cache
	var key = "ClientRepository.GetClients" + strconv.Itoa(int(lastID))
	found, cache := repo.getClientsCache(key)
	if found {
		return cache, nil
	}

	query := `
	select id, email, failedAttempts, lastAttempt, expiredDate, createdOn, updatedOn, enabled
	from client
	where id < $1
	order BY id DESC limit 20
    `
	rows, err := repo.db.Slave.Conn.Query(query, lastID)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	result, err := getClientsFromRows(rows)
	if err != nil {
		return nil, err
	}

	// update cache
	repo.setClientsCache(key, result)
	return result, nil
}

// GetClient method
func (repo *ClientRepository) GetClient(id int64) (models.Client, error) {
	// get data from cache
	var key = "ClientRepository.GetClient" + strconv.Itoa(int(id))
	found, cache := repo.getClientCache(key)
	if found {
		return cache, nil
	}

	var client models.Client
	query := `
	select id, email, failedAttempts, lastAttempt, expiredDate, createdOn, updatedOn, enabled
	from client
	where id=$1 and enabled=true
	`
	row := repo.db.Slave.Conn.QueryRow(query, id)
	err := row.Scan(&client.ID, &client.Email, &client.FailedAttempts, &client.LastAttempt,
		&client.ExpiredDate, &client.CreatedOn, &client.UpdatedOn, &client.Enabled)
	if err != nil {
		return models.Client{}, err
	}

	// update cache
	repo.setClientCache(key, client)
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

// getClientCache method get cache for client
func (repo *ClientRepository) getClientCache(key string) (bool, models.Client) {
	var data models.Client
	if x, found := repo.cache.Get(key); found {
		data = x.(models.Client)
		return found, data
	}
	return false, data
}

// setClientCache method set cache for client
func (repo *ClientRepository) setClientCache(key string, client models.Client) {
	if client != (models.Client{}) {
		clientRepositoryCacheKeys = append(clientRepositoryCacheKeys, key)
		repo.cache.Set(key, client, 5*time.Minute)
	}
}

// getClientCache method get cache for clients
func (repo *ClientRepository) getClientsCache(key string) (bool, []models.Client) {
	var data []models.Client
	if x, found := repo.cache.Get(key); found {
		data = x.([]models.Client)
		return found, data
	}
	return false, data
}

// setClientCache method set cache for clients
func (repo *ClientRepository) setClientsCache(key string, clients []models.Client) {
	if len(clients) > 0 {
		clientRepositoryCacheKeys = append(clientRepositoryCacheKeys, key)
		repo.cache.Set(key, clients, 10*time.Minute)
	}
}

// deleteCache method to delete
func (repo *ClientRepository) deleteCache() {
	if len(clientRepositoryCacheKeys) > 0 {
		for i := 0; i < len(clientRepositoryCacheKeys); i++ {
			repo.cache.Delete(clientRepositoryCacheKeys[i])
		}
		clientRepositoryCacheKeys = []string{}
	}
}
