package repositories

import (
	"context"
	"fmt"

	cache "github.com/greatfocus/gf-sframe/cache"
	"github.com/greatfocus/gf-sframe/database"
	"github.com/greatfocus/gf-user/models"
)

// OtpRepository struct
type OtpRepository struct {
	db    *database.Conn
	cache *cache.Cache
}

// Init method
func (repo *OtpRepository) Init(db *database.Conn, cache *cache.Cache) {
	repo.db = db
	repo.cache = cache
}

// Create method
func (repo *OtpRepository) Create(ctx context.Context, otp models.Otp, channel string) (models.Otp, error) {
	statement := `
    INSERT INTO otp (userId, token, channel, expiredDate)
    VALUES ($1, $2, $3, $4)
    returning id
  `
	var id int64
	err := repo.db.Insert(ctx, statement, otp.UserID, otp.Token, channel, otp.ExpiredDate).Scan(&id)
	if err != nil {
		return otp, err
	}
	createdOtp := otp
	createdOtp.ID = id
	return createdOtp, nil
}

// GetByToken method
func (repo *OtpRepository) GetByToken(ctx context.Context, userID int64, token int64) (models.Otp, error) {
	query := `
	SELECT id, token, expiredDate
	FROM otp
	WHERE userId = $1 AND token = $2 AND verified = false
    `
	row := repo.db.Select(ctx, query, userID, token)
	var otp models.Otp
	err := row.Scan(&otp.ID, &otp.Token, &otp.ExpiredDate)
	if err != nil {
		return models.Otp{}, err
	}

	return otp, nil
}

// Update method
func (repo *OtpRepository) Update(ctx context.Context, otp models.Otp) error {
	query := `
    UPDATE otp
	SET 
		verified=$2,
		updatedOn=CURRENT_TIMESTAMP
    WHERE id=$1
  	`

	res, err := repo.db.Update(ctx, query, otp.ID, true)
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

// Delete method
func (repo *OtpRepository) Delete(ctx context.Context, id int64) error {
	query := `
    delete from otp where id=$1
  	`
	res, err := repo.db.Delete(ctx, query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got updated Otp for %d", id)
	}

	return nil
}
