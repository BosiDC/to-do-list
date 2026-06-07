package repository

import (
	"database/sql"
	"errors"

	"to-do-list/internal/model"
)

type TodosRepository struct {
	db *sql.DB
}

func NewTodosRepository(db *sql.DB) *TodosRepository {
	return &TodosRepository{db: db}
}

func (r *TodosRepository) ListTodos() ([]model.Todo, error) {
	rows, err := r.db.Query("SELECT id, title, completed, created_at FROM todos ORDER BY created_at ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var todo model.Todo
		var completed int
		if err := rows.Scan(&todo.ID, &todo.Title, &completed, &todo.CreatedAt); err != nil {
			return nil, err
		}
		todo.Completed = completed != 0
		todos = append(todos, todo)
	}

	return todos, rows.Err()
}

func (r *TodosRepository) FindByID(id int64) (*model.Todo, error) {
	var todo model.Todo
	var completed int
	row := r.db.QueryRow("SELECT id, title, completed, created_at FROM todos WHERE id = ?", id)
	if err := row.Scan(&todo.ID, &todo.Title, &completed, &todo.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	todo.Completed = completed != 0
	return &todo, nil
}

func (r *TodosRepository) CreateTodo(title string) (*model.Todo, error) {
	result, err := r.db.Exec("INSERT INTO todos (title, completed) VALUES (?, ?)", title, 0)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.FindByID(id)
}

func (r *TodosRepository) UpdateTodo(todo *model.Todo) error {
	_, err := r.db.Exec("UPDATE todos SET title = ?, completed = ? WHERE id = ?", todo.Title, todo.Completed, todo.ID)
	return err
}

func (r *TodosRepository) DeleteTodo(id int64) error {
	_, err := r.db.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}
