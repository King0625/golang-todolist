package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/King0625/golang-todolist/internal/model"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *model.Todo) error
	GetAllByUserId(ctx context.Context, userID int) ([]*model.Todo, error)
	GetById(ctx context.Context, id int) (*model.Todo, error)
	UpdateById(ctx context.Context, id int, title, content string, done bool) error
	MarkDoneById(ctx context.Context, id int) error
	DeleteById(ctx context.Context, id int) error
}

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (t *todoRepository) Create(ctx context.Context, todo *model.Todo) error {
	insertTodoQuery := `INSERT INTO todos (user_id, title, content) VALUES(?,?,?)`

	result, err := t.db.ExecContext(ctx, insertTodoQuery,
		todo.UserID,
		todo.Title,
		todo.Content,
	)

	if err != nil {
		return err
	}

	newId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	todo.ID = int(newId)

	return nil
}

func (t *todoRepository) GetAllByUserId(ctx context.Context, userID int) ([]*model.Todo, error) {
	query := "SELECT * FROM todos WHERE user_id = ?"
	rows, err := t.db.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	var todos []*model.Todo

	for rows.Next() {
		var todo model.Todo

		err = rows.Scan(
			&todo.ID,
			&todo.UserID,
			&todo.Title,
			&todo.Content,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.Done,
		)

		if err != nil {
			return nil, err
		}

		todos = append(todos, &todo)
	}

	return todos, nil
}

func (t *todoRepository) GetById(ctx context.Context, id int) (*model.Todo, error) {
	query := `SELECT * FROM todos WHERE id = ?`
	var todo model.Todo

	err := t.db.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Content,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.Done,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (t *todoRepository) UpdateById(ctx context.Context, id int, title, content string, done bool) error {
	query := "UPDATE todos SET title = ?, content = ?, updatedAt = ?, done = ? WHERE id = ?"
	result, err := t.db.ExecContext(ctx, query, title, content, time.Now(), done, id)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (t *todoRepository) MarkDoneById(ctx context.Context, id int) error {
	query := "UPDATE todos SET done = 1, updatedAt = ? WHERE id = ?"
	result, err := t.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func (t *todoRepository) DeleteById(ctx context.Context, id int) error {
	query := `DELETE FROM todos WHERE id = ?`
	result, err := t.db.Exec(query, id)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
