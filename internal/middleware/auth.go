package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/King0625/golang-todolist/pkg/utils"
)

type contextKey string

const Unauthorized = "UNAUTHORIZED"
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

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var message string
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			message = "Missing or invalid Authorization header"
			utils.RespondError(w, http.StatusUnauthorized, Unauthorized, message, nil)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := utils.ParseJWT(tokenStr)
		if err != nil {
			message = "invalid token"
			utils.RespondError(w, http.StatusUnauthorized, Unauthorized, message, nil)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
