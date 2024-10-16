package repository

import (
	"fmt"

	todogo "github.com/cheboxarov/todo-go"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user todogo.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		logrus.Errorf("error to scan id, after creating user: %s", err.Error())
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username string, password string) (todogo.User, error) {
	var user todogo.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)
	if err != nil {
		logrus.Errorf("error to get ser from database: %s", err.Error())
	}
	return user, err
}
