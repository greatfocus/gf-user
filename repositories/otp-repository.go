package repositories

import (
	"context"
	"errors"

	"github.com/greatfocus/gf-sframe/database"
	"github.com/greatfocus/gf-user/models"
	cache "github.com/patrickmn/go-cache"
)

// OtpRepository struct
type OtpRepository struct {
	db    database.Database
	cache *cache.Cache
}

// Init method
func (repo *OtpRepository) Init(database database.Database, cache *cache.Cache) {
	repo.db = database
	repo.cache = cache
}

// Create method
func (repo *OtpRepository) Create(ctx context.Context, enKey string, otp models.Otp) (models.Otp, error) {
	statement := `
    INSERT INTO otp (token, channel, expiredDate)
    VALUES (PGP_SYM_ENCRYPT($1, '` + enKey + `'), $2, $3)
	returning id
  `
	_, inserted := repo.db.Insert(ctx, statement, otp.Token, otp.Channel, otp.ExpiredDate)
	if !inserted {
		return otp, errors.New("create otp failed")
	}
	createdOtp := otp
	return createdOtp, nil
}

// GetByToken method
func (repo *OtpRepository) GetByToken(ctx context.Context, enKey string, token int64) (models.Otp, error) {
	query := `
	SELECT id, channel, expiredDate, verified, active
	FROM otp
	WHERE pgp_sym_decrypt(token::bytea, '` + enKey + `') = $1 and verified = false and active = true 
    `
	row := repo.db.Select(ctx, query, token)
	var otp models.Otp
	err := row.Scan(&otp.ID, &otp.Channel, &otp.ExpiredDate, &otp.Verified, &otp.Active)
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
		updatedOn=CURRENT_TIMESTAMP,
		verified=$2,
		active=$3
    WHERE id=$1
  	`
	updated := repo.db.Update(ctx, query, otp.ID, otp.Verified, otp.Active)
	if !updated {
		return errors.New("update otp failed")
	}
	return nil
}

// Delete method
func (repo *OtpRepository) Delete(ctx context.Context, id int64) error {
	query := `
    delete from otp where id=$1
  	`
	deleted := repo.db.Delete(ctx, query, id)
	if !deleted {
		return errors.New("update otp failed")
	}
	return nil
}
