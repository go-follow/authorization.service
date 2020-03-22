package transport_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-follow/authorization.service/config"
	"github.com/go-follow/authorization.service/internal/service"
	"github.com/go-follow/authorization.service/internal/service/jwttoken"
	"github.com/go-follow/authorization.service/internal/service/model"
	"github.com/go-follow/authorization.service/internal/service/store/sqlstore"
	"github.com/go-follow/authorization.service/internal/service/transport"
	"github.com/go-follow/authorization.service/pkg/db"
	"github.com/go-follow/authorization.service/pkg/logger"
	"github.com/stretchr/testify/assert"
)

var serv *httptest.Server
var svc service.Service
var s *sqlstore.Store

func TestMain(m *testing.M) {
	router := chi.NewRouter()
	cfg := config.TestConfig()
	conn, err := db.Connect(cfg.Db)
	if err != nil {
		logger.Fatal(err)
	}
	defer conn.Close()
	s = sqlstore.New(conn)
	svc = service.New(s, cfg)
	router.Route("/api/v1", func(rapi chi.Router) {
		transport.NewHTTP(svc, rapi)
	})
	serv = httptest.NewServer(router)
	defer serv.Close()
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	cases := []struct {
		name     string
		wantRes  map[string]string
		wantCode int
	}{
		{
			name:     "success",
			wantRes:  map[string]string{"success": "pong"},
			wantCode: 200,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := serv.Client().Get(fmt.Sprintf("%s/api/v1/ping", serv.URL))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			var v map[string]string
			json.NewDecoder(res.Body).Decode(&v)
			assert.Equal(t, c.wantRes, v)
			assert.Equal(t, c.wantCode, res.StatusCode)
		})
	}
}

func TestCheckAccess(t *testing.T) {
	cases := []struct {
		name      string
		keyCookie string
		token     func() string
		wantRes   map[string]string
		wantCode  int
	}{
		{
			name:      "unknown_key",
			keyCookie: "tok",
			token:     func() string { return "" },
			wantRes:   map[string]string{"error": "invalid authorization token"},
			wantCode:  401,
		},
		{
			name:      "empty_token",
			keyCookie: "token",
			token:     func() string { return "" },
			wantRes:   map[string]string{"error": "invalid authorization token"},
			wantCode:  401,
		},
		{
			name:      "success",
			keyCookie: "token",
			wantRes:   map[string]string{"user_role": "admin"},
			token: func() string {
				cfg := config.TestConfig()
				token, err := jwttoken.Create("admin", "admin", cfg.Server.Secret)
				if err != nil {
					return ""
				}
				return token
			},
			wantCode: 200,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet,
				fmt.Sprintf("%s/api/v1/check_access", serv.URL), nil)
			assert.NoError(t, err)
			req.AddCookie(&http.Cookie{
				Name:     c.keyCookie,
				Value:    c.token(),
				Path:     "/",
				Expires:  time.Now().Add(time.Hour * 3),
				MaxAge:   0,
				HttpOnly: true,
			})
			res, err := serv.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			var v map[string]string
			json.NewDecoder(res.Body).Decode(&v)
			assert.Equal(t, c.wantRes, v)
			assert.Equal(t, c.wantCode, res.StatusCode)
		})
	}
}

func TestLogin(t *testing.T) {
	_, id, err := svc.AddUser(&model.User{
		Name:     "admin",
		Password: "123456",
		Enable:   true,
		Role:     "admin",
	})
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		if err := s.User().Delete(id); err != nil {
			logger.Fatal(err)
		}
	}()
	cases := []struct {
		name     string
		userName string
		password string
		wantCode int
	}{
		{
			name:     "not_found_user",
			userName: "unknown",
			password: "123456",
			wantCode: 401,
		},
		{
			name:     "not_valid_password",
			userName: "admin",
			password: "111",
			wantCode: 401,
		},
		{
			name:     "success",
			userName: "admin",
			password: "123456",
			wantCode: 200,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet,
				fmt.Sprintf("%s/api/v1/login", serv.URL), nil)
			if err != nil {
				t.Fatal(err)
			}
			req.SetBasicAuth(c.userName, c.password)

			res, err := serv.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			assert.Equal(t, c.wantCode, res.StatusCode)
		})
	}
}

func TestAddUser(t *testing.T) {
	cases := []struct {
		name      string
		keyCookie string
		token     func() string
		userName  string
		password  string
		role      string
		wantCode  int
	}{
		{
			name:      "not_cookie",
			keyCookie: "tok",
			token:     func() string { return "" },
			wantCode:  400,
		},
		{
			name:      "not_valid_token",
			keyCookie: "token",
			token:     func() string { return "" },
			wantCode:  401,
		},
		{
			name:      "not_admin",
			keyCookie: "token",
			token: func() string {
				cfg := config.TestConfig()
				token, err := jwttoken.Create("unknown", "editor", cfg.Server.Secret)
				if err != nil {
					return ""
				}
				return token
			},
			wantCode: 403,
		},
		{
			name:      "not_valid_user",
			keyCookie: "token",
			token: func() string {
				cfg := config.TestConfig()
				token, err := jwttoken.Create("unknown", "admin", cfg.Server.Secret)
				if err != nil {
					return ""
				}
				return token
			},
			wantCode: 400,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost,
				fmt.Sprintf("%s/api/v1/add_user", serv.URL), nil)
			if err != nil {
				t.Fatal(err)
			}
			req.AddCookie(&http.Cookie{
				Name:     c.keyCookie,
				Value:    c.token(),
				Path:     "/",
				Expires:  time.Now().Add(time.Hour * 3),
				MaxAge:   0,
				HttpOnly: true,
			})
			res, err := serv.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			assert.Equal(t, c.wantCode, res.StatusCode)
		})
	}
}

func TestLogout(t *testing.T) {
	cases := []struct {
		name     string
		wantCode int
	}{
		{
			name:     "success",
			wantCode: 200,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := serv.Client().Get(fmt.Sprintf("%s/api/v1/logout", serv.URL))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			assert.Equal(t, c.wantCode, res.StatusCode)
		})
	}
}
