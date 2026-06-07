package service

import (
	"to-do-list/internal/model"
	"to-do-list/internal/repository"
)

type TodosService struct {
	repo *repository.TodosRepository
}

func NewTodosService(repo *repository.TodosRepository) *TodosService {
	return &TodosService{repo: repo}
}

func (s *TodosService) ListTodos() ([]model.Todo, error) {
	return s.repo.ListTodos()
}

func (s *TodosService) CreateTodo(title string) (*model.Todo, error) {
	return s.repo.CreateTodo(title)
}

func (s *TodosService) GetTodo(id int64) (*model.Todo, error) {
	return s.repo.FindByID(id)
}

func (s *TodosService) UpdateTodo(todo *model.Todo) error {
	return s.repo.UpdateTodo(todo)
}

func (s *TodosService) DeleteTodo(id int64) error {
	return s.repo.DeleteTodo(id)
}
