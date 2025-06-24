package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("myjwtsecret")

func NewToken(username string, userID int) (accessToken string, err error) {
	accessClaims := jwt.MapClaims{
		"userID":   userID,
		"username": username,
		"exp":      time.Now().Add(2 * time.Hour).Unix(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = at.SignedString(jwtSecret)
	return
}

func ParseJWT(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	return int(claims["userID"].(float64)), nil
}
