package inmemory

import (
	user "birtday_service/models/user"

	"sync"
	"time"
)

var _ user.UserStorage = (*UserStorage)(nil)

type UserStorage struct {
	sync.Mutex

	Users  map[string]user.User
	NextID int
}

func NewUserStorage() *UserStorage {
	users := make(map[string]user.User)
	return &UserStorage{Users: users, NextID: 0}
}

func (s *UserStorage) AddEmp(name string, email string, password string, birthday time.Time) error {
	s.Lock()
	defer s.Unlock()

	user := user.User{ID: s.NextID, Name: name, Password: password, Email: email, Birthday: birthday}
	s.Users[email] = user
	s.NextID++
	return nil
}

func (s *UserStorage) GetByEmail(email string) (user.User, error) {
	s.Lock()
	defer s.Unlock()

	usr, ok := s.Users[email]
	if !ok {
		return user.User{}, user.ErrorNotFound
	}
	return usr, nil
}

func (s *UserStorage) GetAll() ([]user.User, error) {
	s.Lock()
	defer s.Unlock()

	allUsers := make([]user.User, 0)
	for _, usr := range s.Users {
		allUsers = append(allUsers, usr)
	}
	return allUsers, nil
}

func (s *UserStorage) DeleteEmp(email string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.Users, email)
	return nil
}

func (s *UserStorage) GetByID(id int) (user.User, error) {
	s.Lock()
	defer s.Unlock()

	for _, usr := range s.Users {
		if usr.ID == id {
			return usr, nil
		}
	}
	return user.User{}, user.ErrorNotFound
}
