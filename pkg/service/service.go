package service

import (
	todogo "github.com/cheboxarov/todo-go"
	"github.com/cheboxarov/todo-go/pkg/repository"
)

type Authorization interface {
	CreateUser(user todogo.User) (int, error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list todogo.TodoList) (int, error)
	GetAll(userId int) ([]todogo.TodoList, error)
	GetById(userId int, listId int) (todogo.TodoList, error)
	Delete(userId int, listId int) error
	Update(userId int, listId int, input todogo.UpdateListInput) error
}

type TodoItem interface {
	Create(userId int, listId int, item todogo.TodoItem) (int, error)
	GetAll(userId int, listId int) ([]todogo.TodoItem, error)
	GetById(userId int, listId int, itemId int) (todogo.TodoItem, error)
	Delete(userId int, itemId int) error
	Update(userId int, itemId int, input todogo.UpdateItemInput) error
}

type Service struct {
	Authorization
	TodoItem
	TodoList
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
