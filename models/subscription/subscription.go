package models

import (
	"errors"
	"fmt"
	"time"
)

var Location = func() *time.Location {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("Error loading location:", err)
		panic(err)
	}
	return location
}()

type Subscription struct {
	SubscriberID     int
	BirthdayPersonID int
	NotificationTime time.Time
}

type SubStorage interface {
	Sub(subscriberID int, celebratorID int, notifTime time.Time) error
	UnSub(subscriberID int, celebratorID int) error
	GetSubsByID(id int) ([]Subscription, error)
	GetSub(subscriberID, celebratorID int) (Subscription, error)
	List() ([]Subscription, error)
	AddYear(subscriberID int, celebratorID int) error
	UpdateSub(subscriberID int, celebratorID int, notificationTime time.Time) error
}

var (
	ErrorNoRows = errors.New("subscriptions not found")
)
