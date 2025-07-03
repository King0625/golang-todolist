package service

import (
	"context"

	"github.com/King0625/golang-todolist/internal/model"
	"github.com/King0625/golang-todolist/internal/repository"
)

type TodoService interface {
	CreateTodo(ctx context.Context, todo *model.Todo) error
	GetTodosByUserId(ctx context.Context, userID int) ([]*model.Todo, error)
	GetTodoById(ctx context.Context, id int) (*model.Todo, error)
	UpdateTodoById(ctx context.Context, id int, title, content string, done bool) error
	MarkTodoDoneById(ctx context.Context, id int) error
	DeleteTodoById(ctx context.Context, id int) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(r repository.TodoRepository) TodoService {
	return &todoService{r}
}

func (t *todoService) CreateTodo(ctx context.Context, todo *model.Todo) error {
	return t.repo.Create(ctx, todo)
}

func (t *todoService) GetTodosByUserId(ctx context.Context, userID int) ([]*model.Todo, error) {
	return t.repo.GetAllByUserId(ctx, userID)
}

func (t *todoService) GetTodoById(ctx context.Context, id int) (*model.Todo, error) {
	return t.repo.GetById(ctx, id)
}

func (t *todoService) UpdateTodoById(ctx context.Context, id int, title, content string, done bool) error {
	return t.repo.UpdateById(ctx, id, title, content, done)
}

func (t *todoService) MarkTodoDoneById(ctx context.Context, id int) error {
	return t.repo.MarkDoneById(ctx, id)
}

func (t *todoService) DeleteTodoById(ctx context.Context, id int) error {
	return t.repo.DeleteById(ctx, id)
}
