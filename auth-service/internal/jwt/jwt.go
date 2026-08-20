package jwt

import (
	"errors"
	"time"

	"github.com/ErenKarakus1/Notification-Platform/auth-service/internal/model"
	"github.com/golang-jwt/jwt"
)

func GenerateToken(jwtSecret string, user model.User) (string, error) {
	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 8).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", errors.New("couldnt generate token")
	}
	return signedString, nil
}
