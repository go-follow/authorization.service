package sqlstore

import (
	"github.com/go-follow/authorization.service/internal/service/store"
	"github.com/jmoiron/sqlx"
)

// Store - хранилище для DB
type Store struct {
	db             *sqlx.DB
	userRepository *UserRepository
}

//New - инициализация хранилища для store
func New(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

//User - репозиторий для User
func (s *Store) User() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
		}
	}
	return s.userRepository
}