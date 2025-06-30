package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/King0625/golang-todolist/middlewares"
	"github.com/King0625/golang-todolist/models"
	"github.com/King0625/golang-todolist/utils"
)

type RegisterPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginSuccessData struct {
	Token string `json:"token"`
}

func Register() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var res JsonResponse
		var payload RegisterPayload
		err := utils.ReadJSON(w, r, &payload)
		if err != nil {
			log.Fatal(err)
			res.Message = "failed to parse json"
			res.Error = err.Error()
			utils.WriteJSON(w, 500, res)
			return
		}
		currentTime := time.Now()

		// var uesr
		user := models.User{
			Email:     payload.Email,
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Password:  payload.Password,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}

		err = models.Register(user)
		if err != nil {
			res.Message = "register failed"
			res.Error = err.Error()
			utils.WriteJSON(w, 500, res)
			return
		}

		res.Message = "success"
		utils.WriteJSON(w, 201, res)
	}
}

func Login() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var res JsonResponse
		var payload LoginPayload
		err := utils.ReadJSON(w, r, &payload)
		if err != nil {
			log.Fatal(err)
			res.Message = "failed to parse json"
			res.Error = err.Error()
			utils.WriteJSON(w, 500, res)
			return
		}

		user, err := models.Login(payload.Email, payload.Password)
		if err != nil {
			res.Message = "login failed"
			res.Error = err.Error()
			utils.WriteJSON(w, 400, res)
			return
		}

		jwtToken, err := utils.NewToken(user.FirstName+user.LastName, user.ID)
		if err != nil {
			res.Message = "Gen token failed"
			res.Error = err.Error()
			utils.WriteJSON(w, 500, res)
			return
		}

		res.Message = "login successfully"
		res.Data = LoginSuccessData{jwtToken}
		utils.WriteJSON(w, 200, res)
	}
}

func GetUserData() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var res JsonResponse

		userID, ok := middlewares.GetUserID(r)
		if !ok {
			res.Message = "Unauthorized"
			res.Error = "Unauthorized"
			utils.WriteJSON(w, http.StatusUnauthorized, res)
			return
		}

		user, err := models.GetUserDataById(userID)
		if user == nil {
			res.Message = "user not found"
			res.Error = err.Error()
			utils.WriteJSON(w, http.StatusNotFound, res)
			return
		}

		res.Message = "success"
		res.Data = user

		utils.WriteJSON(w, http.StatusOK, res)
	}
}
