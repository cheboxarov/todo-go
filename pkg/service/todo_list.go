package service

import (
	"fmt"

	todogo "github.com/cheboxarov/todo-go"
	"github.com/cheboxarov/todo-go/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{
		repo: repo,
	}
}

func (s *TodoListService) Create(userId int, list todogo.TodoList) (int, error) {
	fmt.Println("Создание в сервисе")
	return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]todogo.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) GetById(userId int, listId int) (todogo.TodoList, error) {
	return s.repo.GetById(userId, listId)
}

func (s *TodoListService) Delete(userId int, listId int) error {
	return s.repo.Delete(userId, listId)
}

func (s *TodoListService) Update(userId int, listId int, input todogo.UpdateListInput) error {
	return s.repo.Update(userId, listId, input)
}
