package repomock

import (
	"errors"

	"github.com/jclumbiarres/gofibertmpl/models"
)

func FindByCredentials(email, password string) (*models.User, error) {
	if email == "test@mail.com" && password == "test12345" {
		return &models.User{
			ID:       1,
			Password: "test12345",
		}, nil
	}

	return nil, errors.New("user not found")
}
