package service

import (
	"context"
	"errors"

	"github.com/King0625/golang-todolist/internal/model"
	"github.com/King0625/golang-todolist/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func comparePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type UserService interface {
	Register(ctx context.Context, user *model.User) error
	Login(ctx context.Context, email, password string) (*model.User, error)
	GetUserDataById(ctx context.Context, id int) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{r}
}

func (u *userService) Register(ctx context.Context, user *model.User) error {
	return u.repo.Create(ctx, user)
}

func (u *userService) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	hashedPassword := user.Password
	pwdMatch := comparePassword(hashedPassword, password)

	if !pwdMatch {
		return nil, errors.New("wrong password")
	}

	return user, nil
}

func (u *userService) GetUserDataById(ctx context.Context, id int) (*model.User, error) {
	return u.repo.GetById(ctx, id)
}
