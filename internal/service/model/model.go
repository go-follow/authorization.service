package model

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
)

//User - данные пользователя
type User struct {
	Name           string
	Password       string
	PasswordEncode string
	Enable         bool
	Role           string
}

//Validate - валидация пользователя
func (u *User) Validate() error {
	if strings.TrimSpace(u.Name) == "" {
		return fmt.Errorf("empty name")
	}
	if len([]rune(u.Name)) > 30 {
		return fmt.Errorf("length field name cannot be more 30 characters")
	}
	if len([]rune(strings.TrimSpace(u.Password))) < 6 {
		return fmt.Errorf("length for password cannot be less 6 characters")
	}
	if strings.TrimSpace(u.PasswordEncode) == "" {
		return fmt.Errorf("empty hash_password")
	}
	if len([]rune(u.PasswordEncode)) > 255 {
		return fmt.Errorf("length field  cannot be more 255 characters")
	}
	if strings.TrimSpace(u.Role) != "admin" &&
		strings.TrimSpace(u.Role) != "editor" &&
		strings.TrimSpace(u.Role) != "viewer" {
		return fmt.Errorf("%s not valid role", u.Role)
	}
	return nil
}

//ComparePassword - сравнениние пароля пользователя с хешированной суммой из БД, смешанной с солью
func (u *User) ComparePassword(secret string) bool {
	mac := hmac.New(sha512.New, []byte(secret))
	_, err := mac.Write([]byte(u.Password))
	if err != nil {
		return false
	}
	hash := hex.EncodeToString(mac.Sum(nil))
	return hash == u.PasswordEncode
}

//SetHashPassword - создание хеша для пароля, смешанного с secret
func (u *User) SetHashPassword(secret string) error {
	if strings.TrimSpace(u.Password) == "" {
		return fmt.Errorf("empty password")
	}
	if strings.TrimSpace(secret) == "" {
		return fmt.Errorf("epmtyp secret")
	}
	mac := hmac.New(sha512.New, []byte(secret))
	_, err := mac.Write([]byte(u.Password))
	if err != nil {
		return err
	}
	u.PasswordEncode = hex.EncodeToString(mac.Sum(nil))
	return nil
}
