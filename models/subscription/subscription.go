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
	Sub(subsID int, celebratorID int, birthday time.Time, notifTime time.Time) error
	UnSub(subsID int, celebratorID int) error
	GetSubsByID(id int) ([]Subscription, error)
	List() ([]Subscription, error)
	UpdateSub(subID int, celebratorID int) error
}

var (
	ErrorEmpty    = errors.New("empty subscriptions")
	ErrorNotFound = errors.New("subscriptions not found")
)
