package sqlstore_test

import (
	"os"
	"testing"

	"github.com/go-follow/authorization.service/internal/service/model"
	"github.com/go-follow/authorization.service/internal/service/store"
	"github.com/go-follow/authorization.service/internal/service/store/sqlstore"
	"github.com/go-follow/authorization.service/pkg/db"
	"github.com/stretchr/testify/assert"
)

var s *sqlstore.Store

func TestMain(m *testing.M) {
	cfg := store.TestConfigDB()
	conn, err := db.Connect(cfg)
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()
	s = sqlstore.New(conn)

	os.Exit(m.Run())
}

func TestFind(t *testing.T) {
	u := model.TestUser(t)
	id, err := s.User().Insert(u)
	assert.NoError(t, err)
	defer s.User().Delete(id)
	cases := []struct {
		name     string
		userName string
		userPass string
		wantUser *model.User
		wantErr  bool
	}{
		{
			name:     "empty_name",
			userName: "",
			userPass: "777",
			wantUser: nil,
			wantErr:  false,
		},
		{
			name:     "empty_pass",
			userName: "admin",
			userPass: "",
			wantUser: &model.User{
				Name:           "admin",
				Password:       "",
				PasswordEncode: "777",
				Enable:         true,
				Role:           "admin",
			},
		},
		{
			name:     "not_found",
			userName: "unknown",
			userPass: "123456",
			wantUser: nil,
			wantErr:  false,
		},
		{
			name:     "success",
			userName: "admin",
			userPass: "123456",
			wantUser: &model.User{
				Name:           "admin",
				Password:       "123456",
				PasswordEncode: "777",
				Enable:         true,
				Role:           "admin",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			u, err := s.User().Find(c.userName, c.userPass)
			assert.Equal(t, c.wantUser, u)
			if c.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInsertUser(t *testing.T) {
	cases := []struct {
		name    string
		u       func() *model.User
		wantErr bool
	}{
		{
			name: "empty_password_hash",
			u: func() *model.User {
				u := model.TestUser(t)
				u.PasswordEncode = ""
				return u
			},
			wantErr: true,
		},
		{
			name: "success",
			u: func() *model.User {
				return model.TestUser(t)
			},
			wantErr: false,
		},
		{
			name: "duplicate",
			u: func() *model.User {
				return model.TestUser(t)
			},
			wantErr: true,
		},
	}
	var idUser int32
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id, err := s.User().Insert(c.u())
			if c.wantErr {
				assert.Error(t, err)
				assert.Equal(t, int32(-1), id)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, -1, id)
				assert.NotEqual(t, 0, id)
				idUser = id
			}
		})
	}
	err := s.User().Delete(idUser)
	assert.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	cases := []struct {
		name    string
		id      int32
		wantErr bool
	}{
		{
			name:    "not_found",
			id:      0,
			wantErr: false,
		},
		{
			name:    "success",
			id:      1,
			wantErr: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := s.User().Delete(c.id)
			if c.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
