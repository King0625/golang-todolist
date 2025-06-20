package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/King0625/golang-todolist/models"
	"github.com/King0625/golang-todolist/utils"
)

type RegisterPayload struct {
	Email string `json:"email"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Password string `json:"password"`
}

type LoginPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`	
}

type JsonResponse struct {
	Message string
	Error string
	Data any
}

func Register() func (w http.ResponseWriter, r *http.Request) {
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

		var user models.User
		user.Email = payload.Email
		user.FirstName = payload.FirstName
		user.LastName = payload.LastName
		user.Password = payload.Password
		user.CreatedAt = currentTime
		user.UpdatedAt = currentTime
		
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

func Login() func (w http.ResponseWriter, r *http.Request) {
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

		_, err = models.Login(payload.Email, payload.Password)
		if err != nil {
			res.Message = "login failed"
			res.Error = err.Error()
			utils.WriteJSON(w, 400, res)
			return	
		}

		res.Message = "login successfully"
		utils.WriteJSON(w, 200, res)
	}
}

