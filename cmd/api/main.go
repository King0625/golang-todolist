package main

import (
	"log"
	"net/http"

	"github.com/King0625/golang-todolist/internal/handler"
	"github.com/King0625/golang-todolist/internal/middleware"
	"github.com/King0625/golang-todolist/internal/model"
	"github.com/King0625/golang-todolist/internal/repository"
	"github.com/King0625/golang-todolist/migration"
)

func main() {
	db := repository.InitDB()
	migration.CreateUserTable(db)
	migration.CreateTodoTable(db)

	model.New(db)

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
