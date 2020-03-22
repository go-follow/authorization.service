package store

import "github.com/go-follow/authorization.service/internal/service/model"

//Store - хранилище
type Store interface {
	User() UserRepository
}

//UserRepository - репозиторий для работы с пользователями
type UserRepository interface {
	Insert(u *model.User) (int32, error)
	Find(name, pass string) (*model.User, error)
	Delete(id int32) error
}