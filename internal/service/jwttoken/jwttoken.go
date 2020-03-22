package jwttoken

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//Create - создание jwt токена
func Create(userName, role, secret string) (string, error) {
	if strings.TrimSpace(userName) == "" {
		return "", fmt.Errorf("empty username")
	}
	if strings.TrimSpace(role) == "" {
		return "", fmt.Errorf("empty role")
	}
	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name": interface{}(userName),
		"user_role": role,
		"exp":       interface{}(time.Now().Add(time.Hour * 3).Unix()),
	})
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

//Check - проверка jwt токен на корректность
func Check(token, secret string) (string, error) {
	tok, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписания %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	role, ok := tok.Claims.(jwt.MapClaims)["user_role"].(string)
	if !ok {
		return "", fmt.Errorf("not valid key in jwt token")
	}
	return role, nil
}
