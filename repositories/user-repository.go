package repositories

import (
	"database/sql"
	"fmt"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-user/models"
)

// UserRepository struct
type UserRepository struct {
	db *database.DB
}

// Init method
func (repo *UserRepository) Init(db *database.DB) {
	repo.db = db
}

// CreateUser method
func (repo *UserRepository) CreateUser(user models.User) (models.User, error) {
	statement := `
    insert into users (type, firstName, middleName, lastName, mobileNumber, email, password, expiredDate, createdBy, status)
    values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    returning id
  `
	var id int64
	err := repo.db.Conn.QueryRow(statement, user.Type, user.FirstName, user.MiddleName, user.LastName,
		user.MobileNumber, user.Email, user.Password, user.ExpiredDate, user.CreatedBy, user.Status).Scan(&id)
	if err != nil {
		return user, err
	}
	createdUser := user
	createdUser.ID = id
	return createdUser, nil
}

// GetPasswordByEmail method
func (repo *UserRepository) GetPasswordByEmail(email string) (models.User, error) {
	var user models.User
	query := `
	select id, email, password, failedAttempts, lastAttempt, lastChange, status, enabled
	from users 
	where email = $1 and deleted=false
    `
	row := repo.db.Conn.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FailedAttempts, &user.LastAttempt, &user.LastChange, &user.Status, &user.Enabled)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetByEmail method
func (repo *UserRepository) GetByEmail(email string) (models.User, error) {
	var user models.User
	query := `
	select id, type, firstName, middleName, lastName, mobileNumber, email, failedAttempts, lastAttempt, lastChange, expiredDate, createdOn, updatedOn, status, enabled
	from users 
	where email = $1 and deleted=false
    `
	row := repo.db.Conn.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Type, &user.FirstName, &user.MiddleName, &user.LastName, &user.MobileNumber,
		&user.Email, &user.FailedAttempts, &user.LastAttempt, &user.LastChange, &user.ExpiredDate, &user.CreatedOn, &user.UpdatedOn, &user.Status, &user.Enabled)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetByEmailOrMobileNumber method
func (repo *UserRepository) GetByEmailOrMobileNumber(email string, mobileNumber string) ([]models.User, error) {
	query := `
	select id, type, firstName, middleName, lastName, mobileNumber, email, failedAttempts, lastAttempt, lastChange, expiredDate, createdOn, updatedOn, status, enabled
	from users 
	where email = $1 or mobileNumber = $2
	`
	rows, err := repo.db.Conn.Query(query, email, mobileNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getUsersFromRows(rows)
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
	res, err := repo.db.Conn.Exec(query, user.ID, user.Status, user.Enabled, user.FailedAttempts, user.ExpiredDate)
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
		lastChange=$4,
		status=$5,
		enabled=$6		
    where id=$1
  	`

	res, err := repo.db.Conn.Exec(query, user.ID, user.LastAttempt, user.FailedAttempts, user.LastChange, user.Status, user.Enabled)
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
func (repo *UserRepository) GetUsers(page int64) ([]models.User, error) {
	query := `
	select id, type, firstName, middleName, lastName, mobileNumber, email, failedAttempts, lastAttempt, lastChange, expiredDate, createdOn, updatedOn, status, enabled
	from users 
	where deleted=false
	order BY createdOn limit 50 OFFSET $1-1
    `
	rows, err := repo.db.Conn.Query(query, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Type, &user.FirstName, &user.MiddleName, &user.LastName, &user.MobileNumber,
			&user.Email, &user.FailedAttempts, &user.LastAttempt, &user.LastChange, &user.ExpiredDate, &user.CreatedOn, &user.UpdatedOn, &user.Status, &user.Enabled)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUser method
func (repo *UserRepository) GetUser(id int64) (models.User, error) {
	var user models.User
	query := `
	select id, type, firstName, middleName, lastName, mobileNumber, email, failedAttempts, lastAttempt, lastChange, expiredDate, createdOn, updatedOn, status, enabled
	from users 
	where id=$1 and deleted=false and enabled=true
	`
	row := repo.db.Conn.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Type, &user.FirstName, &user.MiddleName, &user.LastName, &user.MobileNumber,
		&user.Email, &user.FailedAttempts, &user.LastAttempt, &user.LastChange, &user.ExpiredDate, &user.CreatedOn, &user.UpdatedOn, &user.Status, &user.Enabled)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// prepare users row
func getUsersFromRows(rows *sql.Rows) ([]models.User, error) {
	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Type, &user.FirstName, &user.MiddleName, &user.LastName, &user.MobileNumber, &user.Email,
			&user.FailedAttempts, &user.LastAttempt, &user.LastChange, &user.ExpiredDate, &user.CreatedOn, &user.UpdatedOn, &user.Status, &user.Enabled)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
