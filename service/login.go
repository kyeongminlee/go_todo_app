package service

import (
	"context"
	"fmt"
	"go_todo_app/store"
)

type Login struct {
	DB             store.Queryer
	Repository     UserGetter
	TokenGenerator TokenGenerator
}

func (l *Login) Login(ctx context.Context, name, pw string) (string, error) {
	user, err := l.Repository.GetUser(ctx, l.DB, name)
	if err != nil {
		return "", fmt.Errorf("failed to list: %w", err)
	}
	if err := user.ComparePassword(pw); err != nil {
		return "", fmt.Errorf("wrong password: %w", err)
	}
	jwt, err := l.TokenGenerator.GenerateToken(ctx, *user)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return string(jwt), nil
}
