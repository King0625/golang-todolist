package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/King0625/golang-todolist/utils"
)

type contextKey string

const userIDKey contextKey = "userID"

type JsonResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Data    any    `json:"data"`
}

func GetUserID(r *http.Request) (int, bool) {
	id, ok := r.Context().Value(userIDKey).(int)
	return id, ok
}

func JWTAuth(next http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var res JsonResponse
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			res.Message = "Unauthorized"
			res.Error = "Missing or invalid Authorization header"
			utils.WriteJSON(w, http.StatusUnauthorized, res)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := utils.ParseJWT(tokenStr)
		if err != nil {
			res.Message = "Invalid token"
			res.Error = err.Error()
			utils.WriteJSON(w, http.StatusUnauthorized, res)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
