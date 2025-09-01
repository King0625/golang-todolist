package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/King0625/golang-todolist/internal/middleware"
	"github.com/King0625/golang-todolist/internal/model"
	"github.com/King0625/golang-todolist/internal/service"
	"github.com/King0625/golang-todolist/pkg/utils"
	"github.com/go-playground/validator/v10"
)

type TodoHandler struct {
	service  service.TodoService
	validate *validator.Validate
}

func NewTodoHandler(s service.TodoService) *TodoHandler {
	validate := validator.New()
	return &TodoHandler{s, validate}
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
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
		utils.RespondError(w, http.StatusBadRequest, InvalidJSON, message, nil)
		return
	}

	if err = h.validate.Struct(payload); err != nil {
		message = "validation failed"
		details := make(map[string]string)
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			details[e.Field()] = e.ActualTag()
		}
		utils.RespondError(w, http.StatusBadRequest, ValidationError, message, details)
		return
	}

	todo := model.Todo{
		UserID:  userID,
		Title:   payload.Title,
		Content: payload.Content,
	}

	err = h.service.CreateTodo(r.Context(), &todo)
	if err != nil {
		message = "cannot insert todo into db"
		utils.RespondError(w, http.StatusInternalServerError, InternalError, message, nil)
		return
	}

	message = "create todo successfully"
	utils.RespondSuccess(w, http.StatusCreated, message, nil)

}

func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	var message string
	userID, ok := middleware.GetUserID(r)
	if !ok {
		message = "failed to fetch user identity from parsed jwt token"
		utils.RespondError(w, http.StatusUnauthorized, Unauthorized, message, nil)
		return
	}

	todos, err := h.service.GetTodosByUserId(r.Context(), userID)
	if err != nil {
		message = "cannot get todos from db"
		utils.RespondError(w, http.StatusInternalServerError, InternalError, message, nil)
		return
	}

	message = "fetch todos successfully"
	utils.RespondSuccess(w, http.StatusOK, message, todos)
}

func (h *TodoHandler) GetOneTodoByID(w http.ResponseWriter, r *http.Request) {
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

	todo, err := h.service.GetTodoById(r.Context(), todoID)
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

func (h *TodoHandler) UpdateTodoById(w http.ResponseWriter, r *http.Request) {
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

	var payload UpdateTodoPayload

	err = utils.ReadJSONRequest(w, r, &payload)
	if err != nil {
		log.Fatal(err)
		message = "cannot parse json body"
		utils.RespondError(w, http.StatusBadRequest, InvalidJSON, message, nil)
		return
	}

	if err = h.validate.Struct(payload); err != nil {
		message = "validation failed"
		details := make(map[string]string)
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			details[e.Field()] = e.ActualTag()
		}
		utils.RespondError(w, http.StatusBadRequest, ValidationError, message, details)
		return
	}

	todo, err := h.service.GetTodoById(r.Context(), todoID)
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

	err = h.service.UpdateTodoById(r.Context(), todoID, payload.Title, payload.Title, payload.Done)
	if err != nil {
		message = "failed to update the todo in DB"
		utils.RespondError(w, http.StatusInternalServerError, InternalError, message, nil)
		return
	}

	message = "update the todo successfully"
	utils.RespondSuccess(w, http.StatusOK, message, nil)

}

func (h *TodoHandler) MarkTodoDoneById(w http.ResponseWriter, r *http.Request) {
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

	todo, err := h.service.GetTodoById(r.Context(), todoID)
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

	err = h.service.MarkTodoDoneById(r.Context(), todoID)

	if err != nil {
		message = "failed to mark the todo done in DB"
		utils.RespondError(w, http.StatusInternalServerError, InternalError, message, nil)
		return
	}

	message = "mark the todo done successfully"
	utils.RespondSuccess(w, http.StatusOK, message, nil)
}

func (h *TodoHandler) DeleteTodoById(w http.ResponseWriter, r *http.Request) {
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

	todo, err := h.service.GetTodoById(r.Context(), todoID)
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

	err = h.service.DeleteTodoById(r.Context(), todoID)
	if err != nil {
		message = "cannot delete the todo from DB"
		utils.RespondError(w, http.StatusInternalServerError, InternalError, message, nil)
		return
	}

	message = "delete the todo successfully"
	utils.RespondSuccess(w, http.StatusOK, message, nil)

}
