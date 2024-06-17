package postgres

import (
	model "birtday_service/models/user"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID       int       `db:"id" json:"id"`
	Name     string    `db:"name" json:"name"`
	Email    string    `db:"email" json:"email"`
	Password string    `db:"password" json:"password"`
	Birthday time.Time `db:"birthday" json:"birthday"`
}

var _ model.UserStorage = (*PostgresUsers)(nil)

type PostgresUsers struct {
	db *sqlx.DB
}

func NewPostgresUsers(db *sqlx.DB) *PostgresUsers {
	return &PostgresUsers{db: db}
}

func (pu *PostgresUsers) AddEmp(name string, email string, password string, birthday time.Time) error {
	query := `INSERT INTO users (name, email, password, birthday) VALUES ($1, $2, $3, $4)`
	_, err := pu.db.Exec(query, name, email, password, birthday)
	return err
}

func (pu *PostgresUsers) GetByEmail(email string) (model.User, error) {
	var user User
	query := `SELECT id, name, email, password, birthday FROM users WHERE email = $1`
	err := pu.db.Get(&user, query, email)
	if err == sql.ErrNoRows {
		return model.User{}, model.ErrorNotFound
	}
	if err != nil {
		return model.User{}, err
	}
	return convertToModel(user), nil
}

func (pu *PostgresUsers) GetByID(id int) (model.User, error) {
	var user User
	query := `SELECT id, name, email, password, birthday FROM users WHERE id = $1`
	err := pu.db.Get(&user, query, id)
	if err == sql.ErrNoRows {
		return model.User{}, model.ErrorNotFound
	}
	if err != nil {
		return model.User{}, err
	}
	return convertToModel(user), nil
}

func (pu *PostgresUsers) GetAll() ([]model.User, error) {
	var users []User
	query := `SELECT id, name, email, password, birthday FROM users`
	err := pu.db.Select(&users, query)
	if len(users) == 0 {
		return nil, model.ErrorNotFound
	}
	return convertToModels(users), err
}

func (pu *PostgresUsers) DeleteEmp(email string) error {
	query := `DELETE FROM users WHERE email = $1`
	_, err := pu.db.Exec(query, email)
	return err
}

func convertToModel(user User) model.User {
	return model.User{ID: user.ID, Name: user.Name, Email: user.Email, Password: user.Password, Birthday: user.Birthday}
}

func convertToModels(users []User) []model.User {
	var models []model.User
	for _, user := range users {
		models = append(models, model.User{ID: user.ID, Name: user.Name, Email: user.Email, Password: user.Password, Birthday: user.Birthday})
	}
	return models
}
