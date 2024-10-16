package repository

import (
	todogo "github.com/cheboxarov/todo-go"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user todogo.User) (int, error)
	GetUser(username string, password string) (todogo.User, error)
}

type TodoList interface {
	Create(userId int, list todogo.TodoList) (int, error)
	GetAll(userId int) ([]todogo.TodoList, error)
	GetById(userId int, listId int) (todogo.TodoList, error)
	Delete(userId int, listId int) error
	Update(userId int, listId int, input todogo.UpdateListInput) error
}

type TodoItem interface {
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
	}
}
