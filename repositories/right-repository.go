package repositories

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/greatfocus/gf-sframe/database"
	"github.com/greatfocus/gf-user/models"
	cache "github.com/patrickmn/go-cache"
)

// rightRepositoryCacheKeys array
var rightRepositoryCacheKeys = []string{}

// RightRepository struct
type RightRepository struct {
	db    database.Database
	cache *cache.Cache
}

// Init method
func (repo *RightRepository) Init(database database.Database, cache *cache.Cache) {
	repo.db = database
	repo.cache = cache
}

// CreateDefault method
func (repo *RightRepository) CreateDefault(ctx context.Context, right models.Right) (models.Right, error) {
	statement := `
		INSERT INTO rights (roleId, userId, deleted, enabled)
		SELECT id, $1, false, true FROM role WHERE name='Customer'
		returning id
  `
	var id int64
	id, inserted := repo.db.Insert(ctx, statement, right.UserID)
	if !inserted {
		return right, errors.New("create permissions failed")
	}
	createdRight := right
	createdRight.ID = id
	repo.deleteCache()
	return createdRight, nil
}

// GetRight method
func (repo *RightRepository) GetRight(ctx context.Context, userID int64) (models.Right, error) {
	// get data from cache
	var key = "RightRepository.GetRight" + strconv.Itoa(int(userID))
	found, cache := repo.getRightCache(key)
	if found {
		return cache, nil
	}

	query := `
	select rights.userId, role.id as roleId, role.name as role
	from rights
	inner join role on rights.roleId = role.id
	where rights.userId = $1 and rights.deleted=false and rights.enabled=true
	`
	rows := repo.db.Select(ctx, query, userID)
	result, err := getRightsFromRows(rows)
	if err != nil {
		return models.Right{}, err
	}

	// update cache
	repo.setRightCache(key, result)
	return result, nil
}

// prepare right row
func getRightsFromRows(rows *sql.Row) (models.Right, error) {
	var right models.Right
	err := rows.Scan(&right.UserID, &right.RoleID, &right.Role)
	if err != nil {
		return right, err
	}

	return right, nil
}

// getRightCache method get cache for Right
func (repo *RightRepository) getRightCache(key string) (bool, models.Right) {
	var data models.Right
	if x, found := repo.cache.Get(key); found {
		data = x.(models.Right)
		return found, data
	}
	return false, data
}

// setRightCache method set cache for user
func (repo *RightRepository) setRightCache(key string, right models.Right) {
	if right != (models.Right{}) {
		rightRepositoryCacheKeys = append(rightRepositoryCacheKeys, key)
		repo.cache.Set(key, right, 5*time.Minute)
	}
}

// deleteCache method to delete
func (repo *RightRepository) deleteCache() {
	if len(rightRepositoryCacheKeys) > 0 {
		for i := 0; i < len(rightRepositoryCacheKeys); i++ {
			repo.cache.Delete(rightRepositoryCacheKeys[i])
		}
		rightRepositoryCacheKeys = []string{}
	}
}
