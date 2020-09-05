package repositories

import (
	"database/sql"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-user/models"
)

// RightRepository struct
type RightRepository struct {
	db *database.DB
}

// Init method
func (repo *RightRepository) Init(db *database.DB) {
	repo.db = db
}

// CreateDefault method
func (repo *RightRepository) CreateDefault(right models.Right) (models.Right, error) {
	statement := `
		INSERT INTO rights (roleId, userId, deleted, enabled)
		SELECT id, $1, false, true FROM role WHERE name='Customer'
		returning id
  `
	var id int64
	err := repo.db.Conn.QueryRow(statement, right.UserID).Scan(&id)
	if err != nil {
		return right, err
	}
	createdRight := right
	createdRight.ID = id
	return createdRight, nil
}

// GetRight method
func (repo *RightRepository) GetRight(userID int64) (models.Right, error) {
	query := `
	select rights.userId, role.id as roleId, role.name as role
	from rights
	inner join role on rights.roleId = role.id
	where rights.userId = $1 and rights.deleted=false and rights.enabled=true
	`
	rows := repo.db.Conn.QueryRow(query, userID)
	return getRightsFromRows(rows)
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
