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
	DO $$ 
	DECLARE
		roleId INTEGER := (select id from role where name='customer');
	BEGIN 
		INSERT INTO rights (roleId, userId, status, deleted, enabled)
		VALUES
			(roleId, $1, 'RIGHT.APPROVED', false, true)
		returning id
	END $$;
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

// GetRights method
func (repo *RightRepository) GetRights(userID int64) ([]models.Right, error) {
	query := `
	select rights.userId, role.id as roleId, role.name as role, actions.id as actionId, actions.name as action
	from rights
	inner join role on rights.roleId = role.id
	inner join permission on role.id = permission.roleId
	inner join action on permission.actionId = action.id
	where rights.userId = $1 and rights.status = $2 and rights.deleted=false and rights.enabled=true
	`
	rows, err := repo.db.Conn.Query(query, userID, "RIGHT.APPROVED")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return getRightsFromRows(rows)
}

// prepare right row
func getRightsFromRows(rows *sql.Rows) ([]models.Right, error) {
	rights := []models.Right{}
	for rows.Next() {
		var right models.Right
		err := rows.Scan(&right.UserID, &right.RoleID, &right.Role, &right.ActionID, &right.Action)
		if err != nil {
			return nil, err
		}
		rights = append(rights, right)
	}

	return rights, nil
}
