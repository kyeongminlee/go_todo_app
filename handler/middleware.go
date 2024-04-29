package handler

import (
	"go_todo_app/auth"
	"log"
	"net/http"
)

func AuthMiddleware(j *auth.JWTer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			authedRequest, err := j.FillContext(request)
			if err != nil {
				RespondJSON(request.Context(), writer, ErrResponse{
					Message: "not find auth info",
					Details: []string{err.Error()},
				}, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(writer, authedRequest)
		})
	}
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Print("AdminMiddleware")
		if !auth.IsAdmin(request.Context()) {
			RespondJSON(request.Context(), writer, ErrResponse{
				Message: "not admin",
			}, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
