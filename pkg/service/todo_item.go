package service

import (
	todogo "github.com/cheboxarov/todo-go"
	"github.com/cheboxarov/todo-go/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{
		repo:     repo,
		listRepo: listRepo,
	}
}

func (s *TodoItemService) Create(userId int, listId int, item todogo.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(listId, item)
}

func (s *TodoItemService) GetAll(userId int, listId int) ([]todogo.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) GetById(userId int, listId int, itemId int) (todogo.TodoItem, error) {
	return s.repo.GetById(userId, listId, itemId)
}

func (s *TodoItemService) Delete(userId int, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

func (s *TodoItemService) Update(userId int, itemId int, input todogo.UpdateItemInput) error {
	return s.repo.Update(userId, itemId, input)
}
