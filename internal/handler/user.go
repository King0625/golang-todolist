package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/King0625/golang-todolist/internal/middleware"
	"github.com/King0625/golang-todolist/internal/model"
	"github.com/King0625/golang-todolist/internal/service"
	"github.com/King0625/golang-todolist/pkg/utils"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service  service.UserService
	validate *validator.Validate
}

func NewUserHandler(s service.UserService) *UserHandler {
	validate := validator.New()
	return &UserHandler{s, validate}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var message string

	var payload RegisterPayload
	err := utils.ReadJSONRequest(w, r, &payload)
	if err != nil {
		log.Println(err)
		message = "cannot parse json body"
		utils.RespondError(w, http.StatusBadRequest, InvalidJSON, message, nil)
		return
	}

	if err = h.validate.Struct(payload); err != nil {
		message = "Validation failed"
		errs := err.(validator.ValidationErrors)
		details := make(map[string]string)
		for _, e := range errs {
			details[e.Field()] = e.ActualTag()
		}
		utils.RespondError(w, http.StatusBadRequest, ValidationError, message, details)
		return
	}

	currentTime := time.Now()

	// var uesr
	user := model.User{
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Password:  payload.Password,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	err = h.service.Register(r.Context(), &user)
	if err != nil {
		message := "cannot insert user data into db"
		utils.RespondError(w, http.StatusInternalServerError, InternalError, message, nil)
		return
	}

	message = "register user successfully"
	utils.RespondSuccess(w, http.StatusCreated, message, nil)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var message string
	var payload LoginPayload
	err := utils.ReadJSONRequest(w, r, &payload)
	if err != nil {
		log.Fatal(err)
		message = "cannot parse json body"
		utils.RespondError(w, http.StatusBadRequest, InvalidJSON, message, nil)
		return
	}

	if err = h.validate.Struct(payload); err != nil {
		message = "validation failed"
		errs := err.(validator.ValidationErrors)
		details := make(map[string]string)
		for _, e := range errs {
			details[e.Field()] = e.ActualTag()
		}
		utils.RespondError(w, http.StatusBadRequest, ValidationError, message, details)
		return
	}

	user, err := h.service.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		message = "login failed"
		utils.RespondError(w, http.StatusUnauthorized, Unauthorized, message, nil)
		return
	}

	jwtToken, err := utils.NewToken(user.FirstName+user.LastName, user.ID)
	if err != nil {
		message = "failed to issue jwt token from server"
		utils.RespondError(w, http.StatusInternalServerError, InternalError, message, nil)
		return
	}

	message = "login successfully"
	data := LoginSuccessData{jwtToken}
	utils.RespondSuccess(w, http.StatusOK, message, data)
}

func (h *UserHandler) GetUserData(w http.ResponseWriter, r *http.Request) {
	var message string
	userID, ok := middleware.GetUserID(r)
	if !ok {
		message = "failed to fetch user identity from parsed jwt token"
		utils.RespondError(w, http.StatusUnauthorized, Unauthorized, message, nil)
		return
	}

	user, err := h.service.GetUserDataById(r.Context(), userID)
	if user == nil {
		fmt.Println(err)
		message = "user not found"
		utils.RespondError(w, http.StatusNotFound, UserNotFound, message, nil)
		return
	}

	message = "get user data successfully"
	utils.RespondSuccess(w, http.StatusOK, message, user)
}
