package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/greatfocus/gf-frame/cache"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-user/models"
)

// UserRepository struct
type UserRepository struct {
	db    *database.Conn
	cache *cache.Cache
}

// Init method
func (repo *UserRepository) Init(db *database.Conn, cache *cache.Cache) {
	repo.db = db
	repo.cache = cache
}

// CreateUser method
func (repo *UserRepository) CreateUser(user models.User) (models.User, error) {
	statement := `
    insert into users (type, email, password, expiredDate, status)
    values ($1, $2, $3, $4, $5)
    returning id
  `
	var id int64
	err := repo.db.Master.Conn.QueryRow(statement, user.Type, user.Email, user.Password, user.ExpiredDate, user.Status).Scan(&id)
	if err != nil {
		return user, err
	}
	created := user
	created.ID = id
	return created, nil
}

// GetPasswordByEmail method
func (repo *UserRepository) GetPasswordByEmail(email string) (models.User, error) {
	var user models.User
	query := `
	select id, email, password, failedAttempts, lastAttempt, status, enabled
	from users 
	where email = $1 and deleted=false
    `
	row := repo.db.Slave.Conn.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FailedAttempts, &user.LastAttempt, &user.Status, &user.Enabled)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetByEmail method
func (repo *UserRepository) GetByEmail(email string) (models.User, error) {
	// get data from cache
	var key = "UserRepository.GetByEmail" + string(email)
	found, cache := repo.getUserCache(key)
	if found {
		return cache, nil
	}

	var user models.User
	query := `
	select id, type, email, failedAttempts, lastAttempt, expiredDate, createdOn, updatedOn, status, enabled
	from users 
	where email = $1 and deleted=false
    `
	row := repo.db.Slave.Conn.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Type, &user.Email, &user.FailedAttempts, &user.LastAttempt,
		&user.ExpiredDate, &user.CreatedOn, &user.UpdatedOn, &user.Status, &user.Enabled)
	if err != nil {
		return models.User{}, err
	}

	// update cache
	repo.setUserCache(key, user)
	return user, nil
}

// UpdateUser method
func (repo *UserRepository) UpdateUser(user models.User) error {
	query := `
    update users
	set 
		status=$2, 
		enabled=$3,
		failedAttempts=$4,
		expiredDate=$5,
		updatedOn=CURRENT_TIMESTAMP
    where id=$1 and deleted=false
  	`
	res, err := repo.db.Master.Conn.Exec(query, user.ID, user.Status, user.Enabled, user.FailedAttempts, user.ExpiredDate)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got updated User for %d", user.ID)
	}

	return nil
}

// UpdateLoginAttempt method
func (repo *UserRepository) UpdateLoginAttempt(user models.User) error {
	query := `
    update users
	set 
		lastAttempt=$2, 
		failedAttempts=$3,
		status=$4,
		enabled=$5		
    where id=$1
  	`

	res, err := repo.db.Master.Conn.Exec(query, user.ID, user.LastAttempt, user.FailedAttempts, user.Status, user.Enabled)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got updated User for %d", user.ID)
	}

	return nil
}

// GetUsers method
func (repo *UserRepository) GetUsers(lastID int64) ([]models.User, error) {
	// get data from cache
	var key = "UserRepository.GetUsers" + string(lastID)
	found, cache := repo.getUsersCache(key)
	if found {
		return cache, nil
	}

	query := `
	SELECT id, type, email, failedAttempts, lastAttempt, expiredDate, createdOn, updatedOn, status, enabled
	FROM users 
	WHERE id < $1 and deleted=false
	ORDER BY id DESC limit 20
    `
	rows, err := repo.db.Slave.Conn.Query(query, lastID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Type, &user.Email, &user.FailedAttempts,
			&user.LastAttempt, &user.ExpiredDate, &user.CreatedOn, &user.UpdatedOn, &user.Status, &user.Enabled)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// update cache
	repo.setUsersCache(key, users)
	return users, nil
}

// GetUser method
func (repo *UserRepository) GetUser(id int64) (models.User, error) {
	// get data from cache
	var key = "UserRepository.GetByEmail" + string(id)
	found, cache := repo.getUserCache(key)
	if found {
		return cache, nil
	}

	var user models.User
	query := `
	select id, type, email, failedAttempts, lastAttempt, expiredDate, createdOn, updatedOn, status, enabled
	from users 
	where id=$1 and deleted=false and enabled=true
	`
	row := repo.db.Slave.Conn.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Type, &user.Email, &user.FailedAttempts, &user.LastAttempt,
		&user.ExpiredDate, &user.CreatedOn, &user.UpdatedOn, &user.Status, &user.Enabled)
	if err != nil {
		return models.User{}, err
	}

	// update cache
	repo.setUserCache(key, user)
	return user, nil
}

// Delete method
func (repo *UserRepository) Delete(id int64) error {
	query := `
    delete from users where id=$1
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
		return fmt.Errorf("more than 1 record got updated User for %d", id)
	}

	return nil
}

// prepare users row
func getUsersFromRows(rows *sql.Rows) ([]models.User, error) {
	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Type, &user.Email, &user.FailedAttempts, &user.LastAttempt,
			&user.ExpiredDate, &user.CreatedOn, &user.UpdatedOn, &user.Status, &user.Enabled)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// getUserCache method get cache for user
func (repo *UserRepository) getUserCache(key string) (bool, models.User) {
	var data models.User
	if x, found := repo.cache.Get(key); found {
		data = x.(models.User)
		return found, data
	}
	return false, data
}

// setUserCache method set cache for user
func (repo *UserRepository) setUserCache(key string, user models.User) {
	if user != (models.User{}) {
		repo.cache.Set(key, user, 5*time.Minute)
	}
}

// getUsersCache method get cache for User
func (repo *UserRepository) getUsersCache(key string) (bool, []models.User) {
	var data []models.User
	if x, found := repo.cache.Get(key); found {
		data = x.([]models.User)
		return found, data
	}
	return false, data
}

// setUsersCache method set cache for users
func (repo *UserRepository) setUsersCache(key string, users []models.User) {
	if len(users) > 0 {
		repo.cache.Set(key, users, 5*time.Minute)
	}
}
