package sqlstore

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-follow/authorization.service/internal/service/model"
	"github.com/go-follow/authorization.service/pkg/conv"
	"github.com/go-follow/authorization.service/pkg/localtime"
)

//UserRepository - репозиторий для работы с моделью User
type UserRepository struct {
	store *Store
}

type userDb struct {
	ID           sql.NullInt32  `db:"id"`
	PasswordHash sql.NullString `db:"password_hash"`
	Enable       sql.NullBool   `db:"enable"`
	Role         sql.NullString `db:"role"`
}

//Find - имплементация метода Find для UserRepository
func (r *UserRepository) Find(name, pass string) (*model.User, error) {
	u := &userDb{}
	err := r.store.db.Get(u, `SELECT id, password_hash, enable, role FROM public.user
						  WHERE name = $1`, name)
	//не найдено
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &model.User{
		Name:           name,
		Password:       pass,
		PasswordEncode: conv.ToString(u.PasswordHash),
		Enable:         conv.ToBool(u.Enable),
		Role:           conv.ToString(u.Role),
	}, nil
}

//Insert - имплементация метода Insert для UserRepository
func (r *UserRepository) Insert(u *model.User) (int32, error) {
	f, err := r.Find(u.Name, u.Password)
	if err != nil {
		return -1, err
	}
	if f != nil {
		return -1, fmt.Errorf("Имя пользователя %s уже используется", u.Name)
	}
	var id int32
	err = r.store.db.QueryRow(`INSERT INTO public.user (name, password_hash, enable, role, date_create)
					VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		conv.ToSQLString(u.Name), conv.ToSQLString(u.PasswordEncode),
		u.Enable, conv.ToSQLString(u.Role), localtime.Date(time.Now())).Scan(&id)
	if err != nil {
		return -1, err
	}
	if id <= 0 {
		return -1, fmt.Errorf("empty id after INSERT user")
	}
	return id, nil
}

//Delete - удаление пользователя по id
func (r *UserRepository) Delete(id int32) error {
	_, err := r.store.db.Exec(`DELETE FROM public.user WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
