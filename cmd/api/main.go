package main

import (
	"log"
	"net/http"
	"os"

	"github.com/King0625/golang-todolist/internal/db"
	"github.com/King0625/golang-todolist/internal/dto"
	"github.com/King0625/golang-todolist/internal/handler"
	"github.com/King0625/golang-todolist/internal/middleware"
	"github.com/King0625/golang-todolist/internal/repository"
	"github.com/King0625/golang-todolist/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENV")
	if env != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	dsn := os.Getenv("MYSQL_DSN")

	if err := db.RunMigration(dsn); err != nil {
		log.Fatalf("run migration error: %v", err)
	}

	mysqlInstance, err := db.InitMySQL(dsn)
	if err != nil {
		log.Fatal("Cannot init mysql instance")
	}

	userRepo := repository.NewUserRepository(mysqlInstance)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	todoRepo := repository.NewTodoRepository(mysqlInstance)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)

	r := http.NewServeMux()

	r.HandleFunc("POST /users/register", middleware.ValidationMiddleware[dto.RegisterPayload](userHandler.Register))
	r.HandleFunc("POST /users/login", middleware.ValidationMiddleware[dto.LoginPayload](userHandler.Login))
	r.HandleFunc("GET /users/me", middleware.JWTAuth(userHandler.GetUserData))

	r.HandleFunc("POST /todos", middleware.Chain(todoHandler.CreateTodo,
		middleware.JWTAuth,
		middleware.ValidationMiddleware[dto.CreateTodoPayload],
	))
	r.HandleFunc("GET /todos", middleware.JWTAuth(todoHandler.GetTodos))
	r.HandleFunc("GET /todos/{todoID}", middleware.JWTAuth(todoHandler.GetOneTodoByID))
	r.HandleFunc("PUT /todos/{todoID}", middleware.Chain(todoHandler.UpdateTodoById,
		middleware.JWTAuth,
		middleware.ValidationMiddleware[dto.UpdateTodoPayload],
	))
	r.HandleFunc("PATCH /todos/{todoID}/done", middleware.JWTAuth(todoHandler.MarkTodoDoneById))
	r.HandleFunc("DELETE /todos/{todoID}", middleware.JWTAuth(todoHandler.DeleteTodoById))

	log.Fatal(http.ListenAndServe(":11451", r))
}
