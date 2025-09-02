package middleware

import "net/http"

func Chain(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
