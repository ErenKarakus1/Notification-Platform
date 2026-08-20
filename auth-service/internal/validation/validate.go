package validation

import (
	"errors"
	"net/mail"

	"github.com/ErenKarakus1/Notification-Platform/auth-service/internal/model"
)

func ValidateRegisterRequest(req model.RegisterRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return errors.New("invalid email")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if len(req.Password) > 128 {
		return errors.New("password must be at most 128 characters")
	}
	return nil
}

func ValidateLoginRequest(req model.LoginRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return errors.New("invalid email")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
