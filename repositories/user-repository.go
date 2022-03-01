package repositories

import (
	"context"
	"errors"
	"strconv"

	"time"

	"github.com/greatfocus/gf-sframe/database"
	"github.com/greatfocus/gf-user/models"
	cache "github.com/patrickmn/go-cache"
)

// rightRepositoryCacheKeys array
var userRepositoryCacheKeys = []string{}

// UserRepository struct
type UserRepository struct {
	db    database.Database
	cache *cache.Cache
}

// Init method
func (repo *UserRepository) Init(database database.Database, cache *cache.Cache) {
	repo.db = database
	repo.cache = cache
}

// CreateUser method
func (repo *UserRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	statement := `
    insert into users (type, email, pgp_sym_decrypt(password::bytea,'elders2020'), expiredDate, status)
    values ($1, $2, $3, $4, $5)
    returning id
  `
	var id int64
	id, inserted := repo.db.Insert(ctx, statement, user.Type, user.Email, user.Password, user.ExpiredDate, user.Status)
	if !inserted {
		return user, errors.New("create user failed")
	}
	created := user
	created.ID = id
	repo.deleteCache()
	return created, nil
}

// GetPasswordByEmail method
func (repo *UserRepository) GetPasswordByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	query := `
	select id, email, password, failedAttempts, lastAttempt, successLogins, status, enabled
	from users 
	where email = $1 and deleted=false
    `
	row := repo.db.Select(ctx, query, email)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FailedAttempts, &user.LastAttempt, &user.SuccessLogins, &user.Status, &user.Enabled)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetByEmail method
func (repo *UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	// get data from cache
	var key = "UserRepository.GetByEmail" + email
	found, cache := repo.getUserCache(key)
	if found {
		return cache, nil
	}

	var user models.User
	query := `
	select id, type, email, failedAttempts, lastAttempt, successLogins, expiredDate, createdOn, updatedOn, status, enabled
	from users 
	where email = $1 and deleted=false
    `
	row := repo.db.Select(ctx, query, email)
	err := row.Scan(&user.ID, &user.Type, &user.Email, &user.FailedAttempts, &user.LastAttempt,
		&user.SuccessLogins, &user.ExpiredDate, &user.CreatedOn, &user.UpdatedOn, &user.Status, &user.Enabled)
	if err != nil {
		return models.User{}, err
	}

	// update cache
	repo.setUserCache(key, user)
	return user, nil
}

// UpdateUser method
func (repo *UserRepository) UpdateUser(ctx context.Context, user models.User) error {
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
	updated := repo.db.Update(ctx, query, user.ID, user.Status, user.Enabled, user.FailedAttempts, user.ExpiredDate)
	if updated {
		return errors.New("update user failed")
	}

	repo.deleteCache()
	return nil
}

// UpdateLoginAttempt method
func (repo *UserRepository) UpdateLoginAttempt(ctx context.Context, user models.User) error {
	query := `
    update users
	set 
		lastAttempt=$2,
		successLogins=$3,
		failedAttempts=$4,
		status=$5,
		enabled=$6		
    where id=$1
  	`

	updated := repo.db.Update(ctx, query, user.ID, user.LastAttempt, user.SuccessLogins, user.FailedAttempts, user.Status, user.Enabled)
	if !updated {
		return errors.New("update login attempt failed")
	}

	repo.deleteCache()
	return nil
}

// GetUsers method
func (repo *UserRepository) GetUsers(ctx context.Context, lastID int64) ([]models.User, error) {
	// get data from cache
	var key = "UserRepository.GetUsers" + strconv.Itoa(int(lastID))
	found, cache := repo.getUsersCache(key)
	if found {
		return cache, nil
	}

	query := `
	SELECT id, type, email, failedAttempts, lastAttempt, successLogins, successLogins, expiredDate, createdOn, updatedOn, status, enabled
	FROM users 
	WHERE id < $1 and deleted=false
	ORDER BY id DESC limit 20
    `
	rows, err := repo.db.Query(ctx, query, lastID)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Type, &user.Email, &user.FailedAttempts,
			&user.LastAttempt, &user.SuccessLogins, &user.ExpiredDate, &user.CreatedOn, &user.UpdatedOn, &user.Status, &user.Enabled)
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
func (repo *UserRepository) GetUser(ctx context.Context, id int64) (models.User, error) {
	// get data from cache
	var key = "UserRepository.GetByEmail" + strconv.Itoa(int(id))
	found, cache := repo.getUserCache(key)
	if found {
		return cache, nil
	}

	var user models.User
	query := `
	select id, type, email, failedAttempts, lastAttempt, successLogins, expiredDate, createdOn, updatedOn, status, enabled
	from users 
	where id=$1 and deleted=false and enabled=true
	`
	row := repo.db.Select(ctx, query, id)
	err := row.Scan(&user.ID, &user.Type, &user.Email, &user.FailedAttempts, &user.LastAttempt,
		&user.SuccessLogins, &user.ExpiredDate, &user.CreatedOn, &user.UpdatedOn, &user.Status, &user.Enabled)
	if err != nil {
		return models.User{}, err
	}

	// update cache
	repo.setUserCache(key, user)
	return user, nil
}

// Delete method
func (repo *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `
    delete from users where id=$1
  	`
	deleted := repo.db.Delete(ctx, query, id)
	if !deleted {
		return errors.New("delete user failed")
	}

	repo.deleteCache()
	return nil
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
		userRepositoryCacheKeys = append(userRepositoryCacheKeys, key)
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
		userRepositoryCacheKeys = append(userRepositoryCacheKeys, key)
		repo.cache.Set(key, users, 5*time.Minute)
	}
}

// deleteCache method to delete
func (repo *UserRepository) deleteCache() {
	if len(userRepositoryCacheKeys) > 0 {
		for i := 0; i < len(userRepositoryCacheKeys); i++ {
			repo.cache.Delete(userRepositoryCacheKeys[i])
		}
		userRepositoryCacheKeys = []string{}
	}
}
