package fixture

import (
	"go_todo_app/entity"
	"math/rand"
	"strconv"
	"time"
)

func User(user *entity.User) *entity.User {
	result := &entity.User{
		ID:       entity.UserID(rand.Int()),
		Name:     "kyeongmin" + strconv.Itoa(rand.Int())[:5],
		Password: "password",
		Role:     "admin",
		Created:  time.Now(),
		Modified: time.Now(),
	}
	if user == nil {
		return result
	}

	if user.ID != 0 {
		result.ID = user.ID
	}

	if user.Name != "" {
		result.Name = user.Name
	}

	if user.Password != "" {
		result.Password = user.Password
	}

	if user.Role != "" {
		result.Role = user.Role
	}

	if !user.Created.IsZero() {
		result.Created = user.Created
	}

	if !user.Modified.IsZero() {
		result.Modified = user.Modified
	}

	return result
}
