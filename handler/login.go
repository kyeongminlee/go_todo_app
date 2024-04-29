package handler

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"net/http"
)

type Login struct {
	Service   LoginService
	Validator *validator.Validate
}

func (l *Login) ServedHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	var body struct {
		UserName string `json:"user_name" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		RespondJSON(ctx, writer, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	err := l.Validator.Struct(body)
	if err != nil {
		RespondJSON(ctx, writer, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	jwt, err := l.Service.Login(ctx, body.UserName, body.Password)
	if err != nil {
		RespondJSON(ctx, writer, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	response := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: jwt,
	}

	RespondJSON(ctx, writer, response, http.StatusOK)
}
