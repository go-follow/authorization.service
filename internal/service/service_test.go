package service

import (
	"os"
	"testing"

	"github.com/go-follow/authorization.service/config"
	"github.com/go-follow/authorization.service/internal/service/jwttoken"
	"github.com/go-follow/authorization.service/internal/service/model"
	"github.com/go-follow/authorization.service/internal/service/store/sqlstore"
	"github.com/go-follow/authorization.service/pkg/db"
	"github.com/stretchr/testify/assert"
)

var cfg *config.Config
var svc Service
var s *sqlstore.Store

func TestMain(m *testing.M) {
	cfg = config.TestConfig()
	conn, err := db.Connect(cfg.Db)
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()
	s = sqlstore.New(conn)
	svc = New(s, cfg)
	os.Exit(m.Run())
}

func TestAddUser(t *testing.T) {
	cases := []struct {
		name    string
		u       func() *model.User
		wantErr bool
	}{
		{
			name: "empty_password",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = "   "
				return u
			},
			wantErr: true,
		},
		{
			name: "not_validate",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Role = "fsdaf"
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
	var userID int32
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			_, id, err := svc.AddUser(c.u())
			if c.wantErr {
				assert.Error(t, err)
				assert.Equal(t, int32(-1), id)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, int32(-1), id)
				assert.NotEqual(t, int32(0), id)
				userID = id
			}
		})
	}
	err := s.User().Delete(userID)
	assert.NoError(t, err)
}

func TestLogin(t *testing.T) {
	u := model.TestUser(t)
	err := u.SetHashPassword(cfg.Server.Secret)
	assert.NoError(t, err)
	id, err := s.User().Insert(u)
	assert.NoError(t, err)

	cases := []struct {
		name     string
		userName string
		userPass string
		serv     func() Service
		wantCode int
		wantErr  bool
	}{
		{
			name:     "empty_user_name",
			userName: "",
			userPass: "123456",
			serv:     func() Service { return svc },
			wantCode: 401,
			wantErr:  true,
		},
		{
			name:     "empty_user_pass",
			userName: "admin",
			userPass: "",
			serv:     func() Service { return svc },
			wantCode: 401,
			wantErr:  true,
		},
		{
			name:     "not_found_user_name",
			userName: "unknown",
			userPass: "12345",
			serv:     func() Service { return svc },
			wantCode: 401,
			wantErr:  true,
		},
		{
			name:     "not_valid_pass",
			userName: "admin",
			userPass: "12345",
			serv:     func() Service { return svc },
			wantCode: 401,
			wantErr:  true,
		},
		{
			name:     "empty_sercret",
			userName: "admin",
			userPass: "123456",
			serv: func() Service {
				c := config.TestConfig()
				c.Server.Secret = ""
				return New(s, c)
			},
			wantCode: 500,
			wantErr:  true,
		},
		{
			name:     "success",
			userName: "admin",
			userPass: "123456",
			serv:     func() Service { return svc },
			wantCode: 200,
			wantErr:  false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, token, err := c.serv().Login(c.userName, c.userPass)
			if c.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
	err = s.User().Delete(id)
	assert.NoError(t, err)
}

func TestCheckAccess(t *testing.T) {
	cases := []struct {
		name        string
		createToken func() string
		wantRole    string
		wantErr     bool
	}{
		{
			name: "empty_user",
			createToken: func() string {
				token, err := jwttoken.Create("", "admin", cfg.Server.Secret)
				if err != nil {
					return ""
				}
				return token
			},
			wantRole: "",
			wantErr:  true,
		},
		{
			name: "empty_role",
			createToken: func() string {
				token, err := jwttoken.Create("admin", "", cfg.Server.Secret)
				if err != nil {
					return ""
				}
				return token
			},
			wantRole: "",
			wantErr:  true,
		},
		{
			name: "empty_secret",
			createToken: func() string {
				token, err := jwttoken.Create("admin", "admin", "")
				if err != nil {
					return ""
				}
				return token
			},
			wantRole: "",
			wantErr:  true,
		},
		{
			name: "success",
			createToken: func() string {
				token, err := jwttoken.Create("admin", "admin", cfg.Server.Secret)
				if err != nil {
					return ""
				}
				return token
			},
			wantRole: "admin",
			wantErr:  false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			role, err := svc.CheckAccess(c.createToken())
			assert.Equal(t, c.wantRole, role)
			if c.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
