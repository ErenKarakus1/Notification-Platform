package service

import (
	"errors"

	"github.com/ErenKarakus1/Notification-Platform/auth-service/internal/model"
	"github.com/ErenKarakus1/Notification-Platform/auth-service/internal/password"
	"github.com/google/uuid"
)

func CreateUserFromRequest(req model.RegisterRequest) (model.User, error) {
	hashString, err := password.HashPassword(req.Password)
	if err != nil {
		return model.User{}, errors.New("couldnt hash password")
	}
	return model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashString,
	}, nil
}
