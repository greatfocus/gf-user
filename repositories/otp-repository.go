package repositories

import (
	"fmt"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-user/models"
)

// OtpRepository struct
type OtpRepository struct {
	db *database.DB
}

// Init method
func (repo *OtpRepository) Init(db *database.DB) {
	repo.db = db
}

// Create method
func (repo *OtpRepository) Create(otp models.Otp) (models.Otp, error) {
	statement := `
    insert into otp (userId, token, expiredDate)
    values ($1, $2, $3)
    returning id
  `
	var id int64
	err := repo.db.Conn.QueryRow(statement, otp.UserID, otp.Token, otp.ExpiredDate).Scan(&id)
	if err != nil {
		return otp, err
	}
	createdOtp := otp
	createdOtp.ID = id
	return createdOtp, nil
}

// GetByToken method
func (repo *OtpRepository) GetByToken(userID int64, token int64) (models.Otp, error) {
	query := `
	select id, token, expiredDate
	from otp
	where userId = $1 and token = $2 and verified = false
    `
	row := repo.db.Conn.QueryRow(query, userID, token)
	var otp models.Otp
	err := row.Scan(&otp.ID, &otp.Token, &otp.ExpiredDate)
	if err != nil {
		return models.Otp{}, err
	}

	return otp, nil
}

// Update method
func (repo *OtpRepository) Update(otp models.Otp) error {
	query := `
    update otp
	set 
		verified=$2,
		updatedOn=CURRENT_TIMESTAMP
    where id=$1
  	`

	res, err := repo.db.Conn.Exec(query, otp.ID, true)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got Update Otp for %d", otp.ID)
	}

	return nil
}
