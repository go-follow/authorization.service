package service

import (
	"fmt"

	"github.com/go-follow/authorization.service/config"
	"github.com/go-follow/authorization.service/internal/service/jwttoken"
	"github.com/go-follow/authorization.service/internal/service/model"
	"github.com/go-follow/authorization.service/internal/service/store"
)

//App структура данных для сервиса выгрузок(сервис для веб интерфейса)
type App struct {
	s   store.Store
	cfg *config.Config
}

//New инициализация сервиса unloading
func New(s store.Store, cfg *config.Config) *App {
	return &App{
		s:   s,
		cfg: cfg,
	}
}

//Service - сервис авторизации
type Service interface {
	AddUser(u *model.User) (int, int32, error)
	Login(name, pass string) (int, string, error)
	CheckAccess(token string) (string, error)
}

//AddUser - реализация функции добавления нового пользователя
func (a *App) AddUser(u *model.User) (int, int32, error) {
	//Создание хеша из пароля
	if err := u.SetHashPassword(a.cfg.Server.Secret); err != nil {
		return 400, -1, err
	}
	//Валидация входящего пользователя
	if err := u.Validate(); err != nil {
		return 400, -1, err
	}
	//Добавление пользователя в db
	id, err := a.s.User().Insert(u)
	if err != nil {
		return 500, -1, err
	}
	return 200, id, nil
}

//Login - реализация функции аутентификации
func (a *App) Login(name, pass string) (int, string, error) {
	u, err := a.s.User().Find(name, pass)
	if err != nil {
		return 500, "", err
	}
	if u == nil || !u.Enable || !u.ComparePassword(a.cfg.Server.Secret) {
		return 401, "", fmt.Errorf("the username or password you entered is incorrect")
	}

	token, err := jwttoken.Create(name, u.Role, a.cfg.Server.Secret)
	if err != nil {
		return 500, "", err
	}
	return 200, token, nil
}

//CheckAccess - авторизация пользователя по jwt токену
func (a *App) CheckAccess(token string) (string, error) {
	result, err := jwttoken.Check(token, a.cfg.Server.Secret)
	if err != nil {
		return "", err
	}
	return result, nil
}
