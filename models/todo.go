package models

import (
	"fmt"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	Done      bool      `json:"done,omitempty"`
}

func CreateTodo(newTodo Todo) error {
	insertTodoQuery := `INSERT INTO todos (user_id, title, content) VALUES(?,?,?)`

	result, err := db.Exec(insertTodoQuery,
		newTodo.UserID,
		newTodo.Title,
		newTodo.Content,
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

func GetUserTodosByUserID(userID int) ([]*Todo, error) {
	query := "SELECT * FROM todos WHERE user_id = ?"
	rows, err := db.Query(query, userID)

	if err != nil {
		fmt.Println("GetUserTodosByUserID sql query error")
		return nil, err
	}

	var todos []*Todo

	for rows.Next() {
		var todo Todo

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
			fmt.Println("GetUserTodosByUserID() scan todos error")
			return nil, err
		}

		todos = append(todos, &todo)
	}

	return todos, nil
}

func GetOneUserTodoByID(todoID int) (*Todo, error) {
	query := `SELECT * FROM todos WHERE id = ?`
	var todo Todo

	err := db.QueryRow(query, todoID).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Content,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.Done,
	)

	if err != nil {
		fmt.Println("GetOneUserTodoByUserID scan error")
		return nil, err
	}

	return &todo, nil
}

func UpdateUserTodoById(todoID int, title, content string, done bool) error {
	query := "UPDATE todos SET title = ?, content = ?, updatedAt = ?, done = ? WHERE id = ?"
	result, err := db.Exec(query, title, content, time.Now(), done, todoID)
	if err != nil {
		fmt.Println("UpdateUserTodoById sql query error")
		return err
	}
	res, err := result.RowsAffected()
	if err != nil {
		fmt.Println("UpdateUserTodoById RowsAffected error")
		return err
	}
	fmt.Println(res)
	return nil
}

func MarkUserTodoAsDone(todoID int) error {
	query := "UPDATE todos SET done = 1, updatedAt = ? WHERE id = ?"
	result, err := db.Exec(query, time.Now(), todoID)
	if err != nil {
		fmt.Println("MarkUserTodoAsDone sql query error")
		return err
	}
	res, err := result.RowsAffected()
	if err != nil {
		fmt.Println("MarkUserTodoAsDone RowsAffected error")
		return err
	}
	fmt.Println(res)
	return nil
}

func DeleteUserTodoById(todoID int) error {
	query := `DELETE FROM todos WHERE id = ?`
	result, err := db.Exec(query, todoID)
	if err != nil {
		fmt.Println("DeleteUserTodoById sql query error")
		return err
	}
	res, err := result.RowsAffected()
	if err != nil {
		fmt.Println("DeleteUserTodoById RowsAffected error")
		return err
	}
	fmt.Println(res)
	return nil
}
