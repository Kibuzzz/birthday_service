package postgres

import (
	model "birtday_service/models/subscription"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

var _ model.SubStorage = (*PostgresSubs)(nil)

type PostgresSubs struct {
	db *sqlx.DB
}

type subscription struct {
	SubscriberID     int       `db:"subscriber_id"`
	BirthdayPersonID int       `db:"birthday_person_id"`
	NotificationTime time.Time `db:"notification_time"`
}

func toModel(sub subscription) model.Subscription {
	return model.Subscription{SubscriberID: sub.SubscriberID, BirthdayPersonID: sub.BirthdayPersonID, NotificationTime: sub.NotificationTime}
}

func toModels(subs []subscription) []model.Subscription {
	var models []model.Subscription
	for _, sub := range subs {
		models = append(models, model.Subscription{SubscriberID: sub.SubscriberID, BirthdayPersonID: sub.BirthdayPersonID, NotificationTime: sub.NotificationTime})
	}
	return models
}

func NewPostgresSubs(db *sqlx.DB) *PostgresSubs {
	return &PostgresSubs{db: db}
}

func (ps *PostgresSubs) Sub(subsID int, celebratorID int, notifTime time.Time) error {
	query := `INSERT INTO subscriptions (subscriber_id, birthday_person_id, notification_time) 
	          VALUES ($1, $2, $3)`
	tx := ps.db.MustBegin()
	_, err := tx.Exec(query, subsID, celebratorID, notifTime)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresSubs) UnSub(subsID int, celebratorID int) error {
	query := `DELETE FROM subscriptions WHERE subscriber_id = $1 AND birthday_person_id = $2`
	tx := ps.db.MustBegin()
	_, err := tx.Exec(query, subsID, celebratorID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresSubs) GetSubsByID(subID int) ([]model.Subscription, error) {
	query := `SELECT subscriber_id, birthday_person_id, notification_time 
	          FROM subscriptions 
	          WHERE subscriber_id = $1
			  ORDER BY birthday_person_id`

	var subs []subscription
	err := ps.db.Select(&subs, query, subID)
	if len(subs) == 0 {
		return nil, model.ErrorNoRows
	}
	if err != nil {
		return nil, err
	}
	return toModels(subs), nil
}

func (ps *PostgresSubs) List() ([]model.Subscription, error) {
	query := `SELECT * FROM subscriptions ORDER BY notification_time DESC`
	var subs []subscription
	err := ps.db.Select(&subs, query)
	if err != nil {
		return nil, err
	}
	return toModels(subs), nil
}

func (ps *PostgresSubs) AddYear(subsID int, celebratorID int) error {
	query := `UPDATE subscriptions 
	          SET notification_time = notification_time + interval '1 year'
	          WHERE subscriber_id = $1 AND birthday_person_id = $2`

	tx := ps.db.MustBegin()
	_, err := tx.Exec(query, subsID, celebratorID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresSubs) GetSub(subsID, celebratorID int) (model.Subscription, error) {
	query := `SELECT *
			  FROM subscriptions
	          WHERE subscriber_id = $1 AND birthday_person_id = $2`
	tx := ps.db.MustBegin()
	var sub subscription
	err := tx.Get(&sub, query, subsID, celebratorID)
	if err == sql.ErrNoRows {
		return model.Subscription{}, model.ErrorNoRows
	}
	if err != nil {
		tx.Rollback()
		return model.Subscription{}, err
	}
	err = tx.Commit()
	if err != nil {
		return model.Subscription{}, err
	}
	return toModel(sub), nil
}

func (ps *PostgresSubs) UpdateSub(subsID int, celebratorID int, notificationTime time.Time) error {
	query := `UPDATE subscriptions 
	          SET notification_time = $1
	          WHERE subscriber_id = $2 AND birthday_person_id = $3`
	tx := ps.db.MustBegin()
	_, err := tx.Exec(query, notificationTime, subsID, celebratorID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
