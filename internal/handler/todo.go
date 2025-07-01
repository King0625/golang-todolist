package handler

import (
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
		var res JsonResponse

		userID, ok := middleware.GetUserID(r)
		if !ok {
			res.Message = "Unauthorized"
			res.Error = "Unauthorized"
			utils.WriteJSON(w, http.StatusUnauthorized, res)
			return
		}

		var payload CreateTodoPayload

		err := utils.ReadJSON(w, r, &payload)
		if err != nil {
			log.Fatal(err)
			res.Message = "failed to parse json"
			res.Error = err.Error()
			utils.WriteJSON(w, 500, res)
			return
		}

		todo := model.Todo{
			UserID:  userID,
			Title:   payload.Title,
			Content: payload.Content,
		}

		err = model.CreateTodo(todo)

		if err != nil {
			res.Error = err.Error()
			utils.WriteJSON(w, 500, res)
			return
		}
		res.Message = "success"
		utils.WriteJSON(w, 201, res)
	}
}

func GetTodos() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var res JsonResponse
		userID, ok := middleware.GetUserID(r)
		if !ok {
			res.Message = "Unauthorized"
			res.Error = "Unauthorized"
			utils.WriteJSON(w, http.StatusUnauthorized, res)
			return
		}

		todos, err := model.GetUserTodosByUserID(userID)
		if err != nil {
			res.Error = err.Error()
			utils.WriteJSON(w, 500, res)
			return
		}

		res.Message = "Success"
		res.Data = todos
		utils.WriteJSON(w, 200, res)
	}
}

func GetOneTodoByID() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var res JsonResponse

		userID, ok := middleware.GetUserID(r)
		if !ok {
			res.Message = "Unauthorized"
			res.Error = "Unauthorized"
			utils.WriteJSON(w, http.StatusUnauthorized, res)
			return
		}

		todoIDString := r.PathValue("todoID")
		todoID, err := strconv.Atoi(todoIDString)

		if err != nil {
			res.Error = err.Error()
			utils.WriteJSON(w, 400, res)
			return
		}

		todo, err := model.GetOneUserTodoByID(todoID)

		if todo == nil {
			res.Error = err.Error()
			utils.WriteJSON(w, 404, res)
			return
		}

		userIDInTodo := todo.UserID
		if userID != userIDInTodo {
			res.Error = "Cannot access other user's todo!"
			utils.WriteJSON(w, 403, res)
			return
		}

		res.Message = "success"
		res.Data = todo

		utils.WriteJSON(w, http.StatusOK, res)
	}
}

func UpdateTodoById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var res JsonResponse
		userID, ok := middleware.GetUserID(r)
		if !ok {
			res.Message = "Unauthorized"
			res.Error = "Unauthorized"
			utils.WriteJSON(w, http.StatusUnauthorized, res)
			return
		}

		todoIDString := r.PathValue("todoID")
		todoID, err := strconv.Atoi(todoIDString)

		if err != nil {
			res.Error = err.Error()
			utils.WriteJSON(w, 400, res)
			return
		}

		todo, err := model.GetOneUserTodoByID(todoID)
		if todo == nil {
			res.Error = err.Error()
			utils.WriteJSON(w, 404, res)
			return
		}

		userIDInTodo := todo.UserID
		if userIDInTodo != userID {
			res.Error = "Cannot access other user's todo"
			utils.WriteJSON(w, 403, res)
			return
		}

		var payload UpdateTodoPayload

		err = utils.ReadJSON(w, r, &payload)
		if err != nil {
			log.Fatal(err)
			res.Message = "failed to parse json"
			res.Error = err.Error()
			utils.WriteJSON(w, 500, res)
			return
		}

		err = model.UpdateUserTodoById(todoID, payload.Title, payload.Content, payload.Done)

		if err != nil {
			res.Error = err.Error()
			utils.WriteJSON(w, 500, res)
			return
		}

		res.Message = "Success"
		utils.WriteJSON(w, 200, res)
	}
}

func MarkTodoDoneById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var res JsonResponse
		userID, ok := middleware.GetUserID(r)
		if !ok {
			res.Message = "Unauthorized"
			res.Error = "Unauthorized"
			utils.WriteJSON(w, http.StatusUnauthorized, res)
			return
		}

		todoIDString := r.PathValue("todoID")
		todoID, err := strconv.Atoi(todoIDString)

		if err != nil {
			res.Error = err.Error()
			utils.WriteJSON(w, 400, res)
			return
		}

		todo, err := model.GetOneUserTodoByID(todoID)
		if todo == nil {
			res.Error = err.Error()
			utils.WriteJSON(w, 404, res)
			return
		}

		userIDInTodo := todo.UserID
		if userIDInTodo != userID {
			res.Error = "Cannot access other user's todo"
			utils.WriteJSON(w, 403, res)
			return
		}

		err = model.MarkUserTodoAsDone(todoID)

		if err != nil {
			res.Error = err.Error()
			utils.WriteJSON(w, 500, res)
			return
		}

		res.Message = "Success"
		utils.WriteJSON(w, 200, res)
	}
}

func DeleteTodoById() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var res JsonResponse
		userID, ok := middleware.GetUserID(r)
		if !ok {
			res.Message = "Unauthorized"
			res.Error = "Unauthorized"
			utils.WriteJSON(w, http.StatusUnauthorized, res)
			return
		}
		todoIDString := r.PathValue("todoID")
		todoID, err := strconv.Atoi(todoIDString)

		if err != nil {
			res.Error = err.Error()
			utils.WriteJSON(w, 400, res)
			return
		}

		todo, err := model.GetOneUserTodoByID(todoID)
		if todo == nil {
			res.Error = err.Error()
			utils.WriteJSON(w, 404, res)
			return
		}

		userIDInTodo := todo.UserID
		if userIDInTodo != userID {
			res.Error = "Cannot access other user's todo"
			utils.WriteJSON(w, 403, res)
			return
		}

		err = model.DeleteUserTodoById(todoID)

		if err != nil {
			res.Error = err.Error()
			utils.WriteJSON(w, 500, res)
			return
		}

		res.Message = "Success"
		utils.WriteJSON(w, 204, res)
	}
}
