package repositories

import (
	"database/sql"
	"fmt"

	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-user/models"
	"github.com/lib/pq"
)

// NotifyRepository struct
type NotifyRepository struct {
	db *database.DB
}

// Init method
func (repo *NotifyRepository) Init(db *database.DB) {
	repo.db = db
}

// Create method
func (repo *NotifyRepository) Create(notify models.Notify) (models.Notify, error) {
	statement := `
    insert into notify (templateId, operation, uri, channelId, recipient, createdBy, params, status, sent)
    values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    returning id
  `
	var id int64
	err := repo.db.Conn.QueryRow(statement, notify.TemplateID, notify.Operation, notify.URI, notify.ChannelID, notify.Recipient, notify.UserID, pq.Array(notify.Params), notify.Status, notify.Sent).Scan(&id)
	if err != nil {
		return notify, err
	}
	created := notify
	created.ID = id
	return created, nil
}

// GetNotification method
func (repo *NotifyRepository) GetNotification() ([]models.Notify, error) {
	query := `
	select id, templateId, operation, uri, channelId, recipient, params, status, sent
	from notify
	where sent=false and status='queue'
	order BY createdOn limit 50 OFFSET 1-1
    `
	rows, err := repo.db.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getNotifyFromRows(rows)
}

// Update method
func (repo *NotifyRepository) Update(notify models.Notify) error {
	query := `
    update notify
	set 
		status=$2,
		sent=$3,
		updatedOn=CURRENT_TIMESTAMP
    where id=$1
  	`

	res, err := repo.db.Conn.Exec(query, notify.ID, notify.Status, notify.Sent)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("more than 1 record got Update Otp for %d", notify.ID)
	}

	return nil
}

// prepare users row
func getNotifyFromRows(rows *sql.Rows) ([]models.Notify, error) {
	msgs := []models.Notify{}
	for rows.Next() {
		var msg models.Notify
		err := rows.Scan(&msg.ID, &msg.TemplateID, &msg.Operation, &msg.URI, &msg.ChannelID, &msg.Recipient, pq.Array(&msg.Params), &msg.Status, &msg.Sent)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}
