package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/King0625/golang-todolist/internal/db"
	"github.com/King0625/golang-todolist/internal/handler"
	"github.com/King0625/golang-todolist/internal/middleware"
	"github.com/King0625/golang-todolist/internal/model"
	"github.com/King0625/golang-todolist/migration"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var (
		dbUser = os.Getenv("MYSQL_USERNAME")
		dbPass = os.Getenv("MYSQL_PASSWORD")
		dbHost = os.Getenv("MYSQL_HOST")
		dbPort = os.Getenv("MYSQL_PORT")
		dbName = os.Getenv("MYSQL_DBNAME")
	)

	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)
	mysqlInstance, err := db.InitMySQL(dsn)
	if err != nil {
		log.Fatal("Cannot init mysql instance")
	}

	migration.CreateUserTable(mysqlInstance)
	migration.CreateTodoTable(mysqlInstance)

	model.New(mysqlInstance)

	r := http.NewServeMux()

	r.HandleFunc("POST /users/register", handler.Register())
	r.HandleFunc("POST /users/login", handler.Login())
	r.HandleFunc("GET /users/me", middleware.JWTAuth(handler.GetUserData()))

	r.HandleFunc("POST /todos", middleware.JWTAuth(handler.CreateTodo()))
	r.HandleFunc("GET /todos", middleware.JWTAuth(handler.GetTodos()))
	r.HandleFunc("GET /todos/{todoID}", middleware.JWTAuth(handler.GetOneTodoByID()))
	r.HandleFunc("PUT /todos/{todoID}", middleware.JWTAuth(handler.UpdateTodoById()))
	r.HandleFunc("PATCH /todos/{todoID}/done", middleware.JWTAuth(handler.MarkTodoDoneById()))
	r.HandleFunc("DELETE /todos/{todoID}", middleware.JWTAuth(handler.DeleteTodoById()))

	log.Fatal(http.ListenAndServe(":11451", r))
}
