package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/King0625/golang-todolist/pkg/utils"
	"github.com/go-playground/validator/v10"
)

const (
	requestDataKey  contextKey = "requestData"
	invalidJSON     string     = "INVALID_JSON"
	validationError string     = "VALIDATION_ERROR"
)

var validate = validator.New()

func GetValidatedRequest[T any](r *http.Request) T {
	return r.Context().Value(requestDataKey).(T)
}

func ValidationMiddleware[T any](next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req T
		var message string

		err := utils.ReadJSONRequest(w, r, &req)
		if err != nil {
			log.Println(err)
			message = "cannot parse json body"
			utils.RespondError(w, http.StatusBadRequest, invalidJSON, message, nil)
			return
		}

		if err = validate.Struct(req); err != nil {
			message = "validation failed"
			details := make(map[string]string)
			errs := err.(validator.ValidationErrors)
			for _, e := range errs {
				details[e.Field()] = e.ActualTag()
			}
			utils.RespondError(w, http.StatusBadRequest, validationError, message, details)
			return
		}

		ctx := context.WithValue(r.Context(), requestDataKey, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
