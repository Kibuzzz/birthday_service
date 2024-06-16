package models

import (
	"errors"
	"time"
)

type UserStorage interface {
	AddEmp(name string, email string, password string, birthday time.Time) error
	GetByEmail(email string) (User, error)
	GetByID(id int) (User, error)
	GetAll() ([]User, error)
	DeleteEmp(email string) error
}

var ErrorNotFound = errors.New("user not found")

type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Birthday time.Time `json:"birthday"`
}

func (u *User) Age() int {
	year := time.Now().Year()
	age := year - u.Birthday.Year()
	return age
}
