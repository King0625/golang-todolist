package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName,omitempty"`
	LastName  string    `json:"lastName,omitempty"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var db *sql.DB

func New(dbPool *sql.DB) {
	db = dbPool
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

func comparePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Register(newUser User) error {
	_, err := getUserByEmail(newUser.Email)

	if err == nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := hashPassword(newUser.Password)
	if err != nil {
		return err
	}

	insertUserQuery := `INSERT INTO users (email, firstName, lastName, password, createdAt, updatedAt) VALUES(?,?,?,?,?,?)`

	result, err := db.Exec(insertUserQuery,
		newUser.Email,
		newUser.FirstName,
		newUser.LastName,
		hashedPassword,
		newUser.CreatedAt,
		newUser.UpdatedAt,
	)

	if err != nil {
		return err
	}

	newId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println(newId)
	return nil
}

func Login(email, password string) (*User, error) {
	user, err := getUserByEmail(email)
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

func GetUserDataById(userID int) (*User, error) {
	var user User

	getUserByIdQuery := `
SELECT * FROM users WHERE id = ?`

	err := db.QueryRow(getUserByIdQuery, userID).Scan(
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

func getUserByEmail(email string) (*User, error) {
	var user User

	getUserByEmailQuery := `
SELECT * FROM users WHERE email = ?`

	err := db.QueryRow(getUserByEmailQuery, email).Scan(
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
