package inmemory

import (
	model "birtday_service/models/subscription"
	"sync"
	"time"
)

var _ model.SubStorage = (*SubStore)(nil)

type SubStore struct {
	sync.Mutex
	Subs map[int]map[int]model.Subscription
}

func NewSubStore() *SubStore {
	return &SubStore{
		Subs: make(map[int]map[int]model.Subscription),
	}
}

func (ss *SubStore) Sub(subsID int, celebratorID int, notifTime time.Time) error {
	ss.Lock()
	defer ss.Unlock()

	sub := model.Subscription{SubscriberID: subsID, BirthdayPersonID: celebratorID, NotificationTime: notifTime}
	_, ok := ss.Subs[subsID][celebratorID]
	if !ok {
		ss.Subs[subsID] = make(map[int]model.Subscription)
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

func (ss *SubStore) GetSubsByID(subsID int) ([]model.Subscription, error) {
	ss.Lock()
	defer ss.Unlock()
	subs, ok := ss.Subs[subsID]
	if !ok {
		return []model.Subscription{}, model.ErrorNoRows
	}
	allSubs := []model.Subscription{}
	for _, sub := range subs {
		allSubs = append(allSubs, sub)
	}
	return allSubs, nil
}

func (ss *SubStore) List() ([]model.Subscription, error) {
	list := []model.Subscription{}
	for _, subs := range ss.Subs {
		for _, sub := range subs {
			list = append(list, sub)
		}
	}
	return list, nil
}

func (ss *SubStore) GetSub(subsID, celebratorID int) (model.Subscription, error) {
	ss.Lock()
	defer ss.Unlock()

	sub, ok := ss.Subs[subsID][celebratorID]
	if !ok {
		return model.Subscription{}, model.ErrorNoRows
	}
	return sub, nil
}

func (ss *SubStore) UpdateSub(subsID int, celebratorID int, notifTime time.Time) error {
	ss.Lock()
	defer ss.Unlock()

	sub, ok := ss.Subs[subsID][celebratorID]
	if !ok {
		return model.ErrorNoRows
	}
	sub.NotificationTime = notifTime
	ss.Subs[subsID][celebratorID] = sub
	return nil
}

func (ss *SubStore) AddYear(subID int, celebratorID int) error {
	sub, ok := ss.Subs[subID][celebratorID]
	if !ok {
		return model.ErrorNoRows
	}
	sub.NotificationTime = sub.NotificationTime.AddDate(1, 0, 0)
	ss.Subs[subID][celebratorID] = sub
	return nil
}
