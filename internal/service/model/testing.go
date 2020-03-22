package model

import (
	"math/rand"
	"testing"
	"time"
)

//TestUser - вспомогательная функция для получения пользователя
func TestUser(t *testing.T) *User {
	t.Helper()
	return &User{
		Name:           "admin",
		PasswordEncode: "777",
		Password:       "123456",
		Enable:         true,
		Role:           "admin",
	}
}

//TestGenerateString - вспомогательная функция для получения случайной длинны
func TestGenerateString(t *testing.T, length int) string {
	t.Helper()
	var s []byte
	if length <= 0 {
		return ""
	}
	for i := 0; i < length; i++ {
		const min = 33
		const max = 126
		rand.Seed(time.Now().Unix())
		b := rand.Intn(max-min) + min
		s = append(s, byte(b))
	}
	return string(s)
}
