package repository

import (
	"fmt"
	"strings"

	todogo "github.com/cheboxarov/todo-go"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	fmt.Println("Создания репозитория тодо лист")
	return &TodoListPostgres{
		db: db,
	}
}

func (r *TodoListPostgres) Create(userId int, list todogo.TodoList) (int, error) {
	fmt.Println("создание в репо")
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	if _, err = tx.Exec(createUsersListQuery, userId, id); err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]todogo.TodoList, error) {
	var lists []todogo.TodoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1", todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)
	return lists, err
}

func (r *TodoListPostgres) GetById(userId int, listId int) (todogo.TodoList, error) {
	var list todogo.TodoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2", todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)
	return list, err
}

func (r *TodoListPostgres) Delete(userId int, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2", todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)
	return err
}

func (r *TodoListPostgres) Update(userId int, listId int, input todogo.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}
	var setQuery string
	if len(setValues) == 1 {
		setQuery = setValues[0]
	} else {
		setQuery = strings.Join(setValues, ", ")
	}
	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id = $%d AND ul.user_id = $%d", todoListsTable, setQuery, usersListsTable, argId, argId+1)
	args = append(args, listId, userId)
	logrus.Errorf("update sql query: %s", query)
	_, err := r.db.Exec(query, args...)
	return err
}
