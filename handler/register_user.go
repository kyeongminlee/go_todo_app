package handler

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"go_todo_app/entity"
	"net/http"
)

type RegisterUser struct {
	Service   RegisterUserService
	Validator *validator.Validate
}

func (ru *RegisterUser) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	var b struct {
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role" validate:"required"`
	}

	if err := json.NewDecoder(request.Body).Decode(&b); err != nil {
		RespondJSON(ctx, writer, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	if err := ru.Validator.Struct(b); err != nil {
		RespondJSON(ctx, writer, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	user, err := ru.Service.RegisterUser(ctx, b.Name, b.Password, b.Role)
	if err != nil {
		RespondJSON(ctx, writer, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	response := struct {
		ID entity.UserID `json:"id"`
	}{ID: user.ID}
	RespondJSON(ctx, writer, response, http.StatusOK)
}
