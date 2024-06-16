package inmemory

import (
	subscription "birtday_service/models/subscription"
	"sync"
	"time"
)

type SubStore struct {
	sync.Mutex
	Subs map[int]map[int]subscription.Subscription
}

func NewSubStore() *SubStore {
	return &SubStore{
		Subs: make(map[int]map[int]subscription.Subscription),
	}
}

func (ss *SubStore) Sub(subsID int, celebratorID int, birthday time.Time, notifTime time.Time) error {
	ss.Lock()
	defer ss.Unlock()

	hours := time.Duration(notifTime.Hour())
	minutes := time.Duration(notifTime.Minute())

	thisYearBirthday := time.Date(time.Now().Year(), birthday.Month(), birthday.Day(), birthday.Hour(), birthday.Minute(), 0, 0, subscription.Location)
	notificationTime := thisYearBirthday.Add(-time.Hour * hours)
	notificationTime = notificationTime.Add(-time.Minute * minutes)

	sub := subscription.Subscription{SubscriberID: subsID, BirthdayPersonID: celebratorID, NotificationTime: notificationTime}
	_, ok := ss.Subs[subsID][celebratorID]
	if !ok {
		ss.Subs[subsID] = make(map[int]subscription.Subscription)
	}
	ss.Subs[subsID][celebratorID] = sub
	return nil
}

func (ss *SubStore) UnSub(subsID int, celebratorID int) error {
	ss.Lock()
	defer ss.Unlock()
	delete(ss.Subs[subsID], celebratorID)
	return nil
}

func (ss *SubStore) GetSubsByID(subsID int) ([]subscription.Subscription, error) {
	ss.Lock()
	defer ss.Unlock()
	subs, ok := ss.Subs[subsID]
	if !ok {
		return []subscription.Subscription{}, subscription.ErrorEmpty
	}
	allSubs := []subscription.Subscription{}
	for _, sub := range subs {
		allSubs = append(allSubs, sub)
	}
	return allSubs, nil
}

func (ss *SubStore) List() ([]subscription.Subscription, error) {
	list := []subscription.Subscription{}
	for _, subs := range ss.Subs {
		for _, sub := range subs {
			list = append(list, sub)
		}
	}
	return list, nil
}

func (ss *SubStore) UpdateSub(subID int, celebratorID int) error {
	sub, ok := ss.Subs[subID][celebratorID]
	if !ok {
		return subscription.ErrorNotFound
	}
	sub.NotificationTime = sub.NotificationTime.AddDate(1, 0, 0)
	ss.Subs[subID][celebratorID] = sub
	return nil
}
