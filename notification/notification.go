package notification

import (
	subs "birtday_service/models/subscription"
	usr "birtday_service/models/user"
	"fmt"
)

type Notificator struct {
	Users usr.UserStorage
}

func New(users usr.UserStorage) Notificator {
	return Notificator{Users: users}
}

func (n *Notificator) Notify(sub subs.Subscription) error {
	subscriber, err := n.Users.GetByID(sub.SubscriberID)
	if err != nil {
		return err
	}
	birthdayPerson, err := n.Users.GetByID(sub.BirthdayPersonID)
	if err != nil {
		return err
	}
	// Логика уведомления (пока просто пишет в консоль)
	msg := fmt.Sprintf("У %s день рождения %s. Возраст - %d\n", birthdayPerson.Name, birthdayPerson.Birthday, birthdayPerson.Age())
	fmt.Printf("Отправка уведомления подписчику с id %d\nСообщение: %s", subscriber.ID, msg)
	return nil
}
