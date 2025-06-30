package main

import (
	"log"
	"net/http"

	"github.com/King0625/golang-todolist/handlers"
	"github.com/King0625/golang-todolist/middlewares"
	"github.com/King0625/golang-todolist/models"
)

func main() {
	db := InitDB()
	CreateUserTable(db)
	CreateTodoTable(db)

	models.New(db)

	r := http.NewServeMux()

	r.HandleFunc("POST /users/register", handlers.Register())
	r.HandleFunc("POST /users/login", handlers.Login())
	r.HandleFunc("GET /users/me", middlewares.JWTAuth(handlers.GetUserData()))

	r.HandleFunc("POST /todos", middlewares.JWTAuth(handlers.CreateTodo()))
	r.HandleFunc("GET /todos", middlewares.JWTAuth(handlers.GetTodos()))
	r.HandleFunc("GET /todos/{todoID}", middlewares.JWTAuth(handlers.GetOneTodoByID()))
	r.HandleFunc("PUT /todos/{todoID}", middlewares.JWTAuth(handlers.UpdateTodoById()))
	r.HandleFunc("PATCH /todos/{todoID}/done", middlewares.JWTAuth(handlers.MarkTodoDoneById()))
	r.HandleFunc("DELETE /todos/{todoID}", middlewares.JWTAuth(handlers.DeleteTodoById()))

	log.Fatal(http.ListenAndServe(":11451", r))
}
