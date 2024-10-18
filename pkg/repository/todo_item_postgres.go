package repository

import (
	"fmt"

	todogo "github.com/cheboxarov/todo-go"
	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{
		db: db,
	}
}

func (r *TodoItemPostgres) Create(listId int, item todogo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int

	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)
	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListsItemQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	if _, err := tx.Exec(createListsItemQuery, listId, id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userId int, listId int) ([]todogo.TodoItem, error) {
	var items []todogo.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti 
	INNER JOIN %s li on li.item_id = ti.id 
	INNER JOIN %s ul on ul.list_id = li.list_id
	WHERE li.list_id = $1 AND ul.user_id = $2`, todoItemsTable, listsItemsTable, usersListsTable)
	err := r.db.Select(&items, query, listId, userId)
	return items, err
}

func (r *TodoItemPostgres) GetById(userId int, listId int, itemId int) (todogo.TodoItem, error) {
	var item todogo.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti 
	INNER JOIN %s li on li.item_id = ti.id 
	INNER JOIN %s ul on ul.list_id = li.list_id
	WHERE li.list_id = $1 AND ul.user_id = $2 AND li.item_id = $3`, todoItemsTable, listsItemsTable, usersListsTable)
	err := r.db.Get(&item, query, listId, userId, itemId)
	return item, err
}

func (r *TodoItemPostgres) Delete(userId int, itemId int) error {
	query := fmt.Sprintf("DELETE FROM %s ti USING %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2", todoItemsTable, listsItemsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, itemId)
	return err
}

func (r *TodoItemPostgres) Update(userId int, itemId int, input todogo.UpdateItemInput) error {
	qg := NewSetQueryGenerator()
	if input.Title != nil {
		qg.Add("title", *input.Title)
	}
	if input.Description != nil {
		qg.Add("description", *input.Description)
	}

	if input.Done != nil {
		qg.Add("done", *input.Done)
	}

	setQuery, _ := qg.GetSetQuery()
	argId := qg.ArgId
	args := qg.Args
	query := fmt.Sprintf("UPDATE %s ti SET %s FROM %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d", todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)
	args = append(args, userId, itemId)
	_, err := r.db.Exec(query, args...)
	return err
}
