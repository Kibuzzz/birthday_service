package cron

import (
	subscription "birtday_service/models/subscription"
	"birtday_service/notification"
	"log"
	"time"
)

type Cron struct {
	Subs        subscription.SubStorage
	Notificator notification.Notificator
}

func New(subs subscription.SubStorage, notificator notification.Notificator) Cron {
	return Cron{Subs: subs, Notificator: notificator}
}

func (c *Cron) checkBirthday() error {
	subs, err := c.Subs.List()
	if err != nil {
		return err
	}
	now := time.Now()
	for _, sub := range subs {
		if now.After(sub.NotificationTime) {
			if err := c.Notificator.Notify(sub); err != nil {
				return err
			}
			err = c.Subs.UpdateSub(sub.SubscriberID, sub.BirthdayPersonID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Cron) Start() {
	ticker := time.NewTicker(time.Second * 10)
	for range ticker.C {
		err := c.checkBirthday()
		if err != nil {
			log.Print(err)
		}
	}
}
