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

	log.Fatal(http.ListenAndServe(":11451", r))
}
