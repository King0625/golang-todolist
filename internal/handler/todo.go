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

type CreateTodoPayload struct {
	Title   string `json:"title" validate:"required,max=666" example:"sleep"`
	Content string `json:"content" validate:"required,max=6666" example:"sleep forever"`
}

type UpdateTodoPayload struct {
	CreateTodoPayload
	Done bool `json:"done" validate:"required" example:"true"`
}

type TodoHandler struct {
	service  service.TodoService
	validate *validator.Validate
}

func NewTodoHandler(s service.TodoService) *TodoHandler {
	validate := validator.New()
	return &TodoHandler{s, validate}
}

// CreateTodo godoc
// @Summary      Create a new todo
// @Description  Create a todo item with title and content fields
// @Tags         todos
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        todo  body  CreateTodoPayload  true  "Todo to create"
// @Success      201   {object}  utils.SuccessResponse
// @Failure      401   {object}  utils.ErrorResponse
// @Failure      400   {object}  utils.ErrorResponse
// @Router       /todos [post]
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

// GetTodos godoc
// @Summary      Get all todos from authorized user
// @Description  Get all todos from authorized user
// @Tags         todos
// @Security     BearerAuth
// @Produce      json
// @Success      200   {object}  utils.SuccessResponse
// @Failure      401   {object}  utils.ErrorResponse
// @Router       /todos [get]
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

// GetOneTodoByID godoc
// @Summary      get todo by id
// @Description  get one todo via id param
// @Tags         todos
// @Security     BearerAuth
// @Produce      json
// @Param        todoID   path      int  true  "todo ID"
// @Success      200   {object}  utils.SuccessResponse
// @Failure      400   {object}  utils.ErrorResponse
// @Failure      401   {object}  utils.ErrorResponse
// @Failure      403   {object}  utils.ErrorResponse
// @Router       /todos/{todoID} [get]
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

// UpdateTodoById godoc
// @Summary      update todo by id
// @Description  update one todo
// @Tags         todos
// @Security     BearerAuth
// @Produce      json
// @Param        todoID   path      int  true  "todo ID"
// @Param        todo  body  UpdateTodoPayload  true  "Todo to update"
// @Success      200   {object}  utils.SuccessResponse
// @Failure      400   {object}  utils.ErrorResponse
// @Failure      401   {object}  utils.ErrorResponse
// @Failure      403   {object}  utils.ErrorResponse
// @Router       /todos/{todoID} [put]
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

// MarkTodoDoneById godoc
// @Summary      mark todo done by id
// @Description  mark todo done by id
// @Tags         todos
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        todoID   path      int  true  "todo ID"
// @Success      200   {object}  utils.SuccessResponse
// @Failure      400   {object}  utils.ErrorResponse
// @Failure      401   {object}  utils.ErrorResponse
// @Failure      403   {object}  utils.ErrorResponse
// @Router       /todos/{todoID}/done [patch]
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

// DeleteTodoById godoc
// @Summary      delete todo by id
// @Description  delete todo by id
// @Tags         todos
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        todoID   path      int  true  "todo ID"
// @Success      200   {object}  utils.SuccessResponse
// @Failure      400   {object}  utils.ErrorResponse
// @Failure      401   {object}  utils.ErrorResponse
// @Failure      403   {object}  utils.ErrorResponse
// @Router       /todos/{todoID} [delete]
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
