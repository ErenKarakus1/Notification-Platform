package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("couldnt hash password")
	}
	hashString := string(hashedBytes)
	return hashString, nil
}

func CompareHashAndPassword(password_hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
}
