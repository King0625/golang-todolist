package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/King0625/golang-todolist/internal/middleware"
	"github.com/King0625/golang-todolist/internal/model"
	"github.com/King0625/golang-todolist/pkg/utils"
)

type CreateTodoPayload struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateTodoPayload struct {
	CreateTodoPayload
	Done bool `json:"done"`
}

type UpdateTodoStatusPayload struct {
	Done bool `json:"done"`
}

func CreateTodo() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var message string
		userID, ok := middleware.GetUserID(r)
		if !ok {
			message = "failed to fetch user identity from parsed jwt token"
			utils.RespondError(w, http.StatusUnauthorized, Unauthorized, message, nil)
			return
		}

		var payload CreateTodoPayload

		err := utils.ReadJSONRequest(w, r, &payload)
		if err != nil {
			log.Fatal(err)
			message = "cannot parse json body"
			utils.RespondError(w, http.StatusInternalServerError, InvalidJSON, message, nil)
			return
		}

		todo := model.Todo{
			UserID:  userID,
			Title:   payload.Title,
			Content: payload.Content,
		}

		err = model.CreateTodo(todo)

		if err != nil {
			message = "cannot insert todo into db"
			utils.RespondError(w, http.StatusInternalServerError, InternalError, message, nil)
			return
		}

		message = "create todo successfully"
		utils.RespondSuccess(w, http.StatusCreated, message, nil)
	}
}

func GetTodos() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var message string
		userID, ok := middleware.GetUserID(r)
		if !ok {
			message = "failed to fetch user identity from parsed jwt token"
			utils.RespondError(w, http.StatusUnauthorized, Unauthorized, message, nil)
			return
		}

		todos, err := model.GetUserTodosByUserID(userID)
		if err != nil {
			message = "cannot get todos from db"
			utils.RespondError(w, http.StatusInternalServerError, InternalError, message, nil)
			return
		}

		message = "fetch todos successfully"
		utils.RespondSuccess(w, http.StatusOK, message, todos)
	}
}

func GetOneTodoByID() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var message string
		userID, ok := middleware.GetUserID(r)
		if !ok {
			message = "failed to fetch user identity from parsed jwt token"
			utils.RespondError(w, http.StatusUnauthorized, Unauthorized, message, nil)
			return
		}

		todoIDString := r.PathValue("todoID")
		todoID, err := strconv.Atoi(todoIDString)

		if err != nil {
			message = "invalid todoID"
			utils.RespondError(w, http.StatusBadRequest, ValidationError, message, nil)
			return
		}

		todo, err := model.GetOneUserTodoByID(todoID)

		if todo == nil {
			fmt.Println(err)
			message = "todo not found"
			utils.RespondError(w, http.StatusNotFound, TodoNotFound, message, nil)
			return
		}

		userIDInTodo := todo.UserID
		if userID != userIDInTodo {
			message = "this is not your todo"
			utils.RespondError(w, http.StatusForbidden, PermissionDenied, message, nil)
			return
		}

		message = "fetch a todo successfully"
		utils.RespondSuccess(w, http.StatusOK, message, todo)
	}
}

func UpdateTodoById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var message string
		userID, ok := middleware.GetUserID(r)
		if !ok {
			message = "failed to fetch user identity from parsed jwt token"
			utils.RespondError(w, http.StatusUnauthorized, Unauthorized, message, nil)
			return
		}

		todoIDString := r.PathValue("todoID")
		todoID, err := strconv.Atoi(todoIDString)

		if err != nil {
			message = "invalid todoID"
			utils.RespondError(w, http.StatusBadRequest, ValidationError, message, nil)
			return
		}

		todo, err := model.GetOneUserTodoByID(todoID)
		if todo == nil {
			fmt.Println(err)
			message = "todo not found"
			utils.RespondError(w, http.StatusNotFound, TodoNotFound, message, nil)
			return
		}

		userIDInTodo := todo.UserID
		if userIDInTodo != userID {
			message = "this is not your todo"
			utils.RespondError(w, http.StatusForbidden, PermissionDenied, message, nil)
			return
		}

		var payload UpdateTodoPayload

		err = utils.ReadJSONRequest(w, r, &payload)
		if err != nil {
			log.Fatal(err)
			message = "cannot parse json body"
			utils.RespondError(w, http.StatusInternalServerError, InvalidJSON, message, nil)
			return
		}

		err = model.UpdateUserTodoById(todoID, payload.Title, payload.Content, payload.Done)

		if err != nil {
			message = "failed to update the todo in DB"
			utils.RespondError(w, http.StatusInternalServerError, InternalError, message, nil)
			return
		}

		message = "update the todo successfully"
		utils.RespondSuccess(w, http.StatusOK, message, nil)
	}
}

func MarkTodoDoneById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var message string
		userID, ok := middleware.GetUserID(r)
		if !ok {
			message = "failed to fetch user identity from parsed jwt token"
			utils.RespondError(w, http.StatusUnauthorized, Unauthorized, message, nil)
			return
		}

		todoIDString := r.PathValue("todoID")
		todoID, err := strconv.Atoi(todoIDString)

		if err != nil {
			message = "invalid todoID"
			utils.RespondError(w, http.StatusBadRequest, ValidationError, message, nil)
			return
		}

		todo, err := model.GetOneUserTodoByID(todoID)
		if todo == nil {
			fmt.Println(err)
			message = "todo not found"
			utils.RespondError(w, http.StatusNotFound, TodoNotFound, message, nil)
			return
		}

		userIDInTodo := todo.UserID
		if userIDInTodo != userID {
			message = "this is not your todo"
			utils.RespondError(w, http.StatusForbidden, PermissionDenied, message, nil)
			return
		}

		err = model.MarkUserTodoAsDone(todoID)

		if err != nil {
			message = "failed to mark the todo done in DB"
			utils.RespondError(w, http.StatusInternalServerError, InternalError, message, nil)
			return
		}

		message = "mark the todo done successfully"
		utils.RespondSuccess(w, http.StatusOK, message, nil)
	}
}

func DeleteTodoById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var message string
		userID, ok := middleware.GetUserID(r)
		if !ok {
			message = "failed to fetch user identity from parsed jwt token"
			utils.RespondError(w, http.StatusUnauthorized, Unauthorized, message, nil)
			return
		}

		todoIDString := r.PathValue("todoID")
		todoID, err := strconv.Atoi(todoIDString)

		if err != nil {
			message = "invalid todoID"
			utils.RespondError(w, http.StatusBadRequest, ValidationError, message, nil)
			return
		}

		todo, err := model.GetOneUserTodoByID(todoID)
		if todo == nil {
			fmt.Println(err)
			message = "todo not found"
			utils.RespondError(w, http.StatusNotFound, TodoNotFound, message, nil)
			return
		}

		userIDInTodo := todo.UserID
		if userIDInTodo != userID {
			message = "this is not your todo"
			utils.RespondError(w, http.StatusForbidden, PermissionDenied, message, nil)
			return
		}

		err = model.DeleteUserTodoById(todoID)
		if err != nil {
			message = "cannot delete the todo from DB"
			utils.RespondError(w, http.StatusInternalServerError, InternalError, message, nil)
			return
		}

		message = "delete the todo successfully"
		utils.RespondSuccess(w, http.StatusOK, message, nil)
	}
}
