package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/King0625/golang-todolist/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetById(ctx context.Context, id int) (*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *model.User) error {
	_, err := r.GetByEmail(ctx, u.Email)

	if err == nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := hashPassword(u.Password)
	if err != nil {
		return err
	}

	insertUserQuery := `INSERT INTO users (email, firstName, lastName, password, createdAt, updatedAt) VALUES(?,?,?,?,?,?)`

	result, err := r.db.ExecContext(ctx, insertUserQuery,
		u.Email,
		u.FirstName,
		u.LastName,
		hashedPassword,
		u.CreatedAt,
		u.UpdatedAt,
	)

	if err != nil {
		return err
	}

	newId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = int(newId)
	return nil

}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	getUserByEmailQuery := `
SELECT * FROM users WHERE email = ?`

	err := r.db.QueryRowContext(ctx, getUserByEmailQuery, email).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetById(ctx context.Context, id int) (*model.User, error) {
	var user model.User

	getUserByIdQuery := `
SELECT * FROM users WHERE id = ?`

	err := r.db.QueryRowContext(ctx, getUserByIdQuery, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
