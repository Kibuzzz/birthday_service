package postgres

import (
	model "birtday_service/models/subscription"
	"time"

	"github.com/jmoiron/sqlx"
)

type PostgresSubs struct {
	db *sqlx.DB
}

type subscription struct {
	SubscriberID     int       `db:"subscriber_id"`
	BirthdayPersonID int       `db:"birthday_person_id"`
	NotificationTime time.Time `db:"notification_time"`
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

func (pu *PostgresSubs) Sub(subsID int, celebratorID int, birthday time.Time, notifTime time.Time) error {
	hours := time.Duration(notifTime.Hour())
	minutes := time.Duration(notifTime.Minute())

	thisYearBirthday := time.Date(time.Now().Year(), birthday.Month(), birthday.Day(), birthday.Hour(), birthday.Minute(), 0, 0, model.Location)
	notificationTime := thisYearBirthday.Add(-time.Hour * hours)
	notificationTime = notificationTime.Add(-time.Minute * minutes)

	query := `INSERT INTO subscriptions (subscriber_id, birthday_person_id, notification_time) 
	          VALUES ($1, $2, $3)`
	tx := pu.db.MustBegin()
	tx.MustExec(query, subsID, celebratorID, notificationTime)
	_, err := tx.Exec(query, subsID, celebratorID, notificationTime)
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

func (pu *PostgresSubs) UnSub(subsID int, celebratorID int) error {
	query := `DELETE FROM subscriptions WHERE subscriber_id = $1 AND birthday_person_id = $2`
	tx := pu.db.MustBegin()
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

func (pu *PostgresSubs) GetSubsByID(subID int) ([]model.Subscription, error) {
	query := `SELECT subscriber_id, birthday_person_id, notification_time 
	          FROM subscriptions 
	          WHERE subscriber_id = $1`

	var subs []subscription
	err := pu.db.Select(&subs, query, subID)
	if err != nil {
		return nil, err
	}
	if len(subs) == 0 {
		return nil, model.ErrorEmpty
	}
	return toModels(subs), nil
}

func (pu *PostgresSubs) List() ([]model.Subscription, error) {
	query := `SELECT * FROM subscriptions`
	var subs []subscription
	err := pu.db.Select(&subs, query)
	if err != nil {
		return nil, err
	}
	return toModels(subs), nil
}

func (pu *PostgresSubs) UpdateSub(subsID int, celebratorID int) error {
	query := `UPDATE subscriptions 
	          SET notification_time = notification_time + interval '1 year'
	          WHERE subscriber_id = $1 AND birthday_person_id = $2`

	tx := pu.db.MustBegin()
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
