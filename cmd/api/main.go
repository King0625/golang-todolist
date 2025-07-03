package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/King0625/golang-todolist/internal/db"
	"github.com/King0625/golang-todolist/internal/handler"
	"github.com/King0625/golang-todolist/internal/middleware"
	"github.com/King0625/golang-todolist/internal/repository"
	"github.com/King0625/golang-todolist/internal/service"
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
	userRepo := repository.NewUserRepository(mysqlInstance)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	todoRepo := repository.NewTodoRepository(mysqlInstance)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)

	r := http.NewServeMux()

	r.HandleFunc("POST /users/register", userHandler.Register)
	r.HandleFunc("POST /users/login", userHandler.Login)
	r.HandleFunc("GET /users/me", middleware.JWTAuth(userHandler.GetUserData))

	r.HandleFunc("POST /todos", middleware.JWTAuth(todoHandler.CreateTodo))
	r.HandleFunc("GET /todos", middleware.JWTAuth(todoHandler.GetTodos))
	r.HandleFunc("GET /todos/{todoID}", middleware.JWTAuth(todoHandler.GetOneTodoByID))
	r.HandleFunc("PUT /todos/{todoID}", middleware.JWTAuth(todoHandler.UpdateTodoById))
	r.HandleFunc("PATCH /todos/{todoID}/done", middleware.JWTAuth(todoHandler.MarkTodoDoneById))
	r.HandleFunc("DELETE /todos/{todoID}", middleware.JWTAuth(todoHandler.DeleteTodoById))

	log.Fatal(http.ListenAndServe(":11451", r))
}
