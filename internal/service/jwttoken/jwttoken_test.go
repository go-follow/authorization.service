package jwttoken_test

import (
	"testing"

	"github.com/go-follow/authorization.service/internal/service/jwttoken"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name     string
		userName string
		role     string
		secret   string
		wantErr  bool
	}{
		{
			name:     "empty_user_name",
			userName: "",
			role:     "admin",
			secret:   "secret",
			wantErr:  true,
		},
		{
			name:     "empty_role",
			userName: "admin",
			role:     "",
			secret:   "secret",
			wantErr:  true,
		},
		{
			name:     "success",
			userName: "admin",
			role:     "admin",
			secret:   "secret",
			wantErr:  false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			token, err := jwttoken.Create(c.userName, c.role, c.secret)
			if c.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	cases := []struct {
		name     string
		token    func() string
		secret   string
		wantRole string
		wantErr  bool
	}{
		{
			name:     "empty_token",
			token:    func() string { return "" },
			secret:   "secret",
			wantRole: "",
			wantErr:  true,
		},
		{
			name: "not_valid_secret",
			token: func() string {
				token, err := jwttoken.Create("admin", "admin", "secret")
				if err != nil {
					return ""
				}
				return token
			},
			secret: "123",
			wantRole: "",
			wantErr: true,
		},
		{
			name: "success",
			token: func() string {
				token, err := jwttoken.Create("admin", "admin", "secret")
				if err != nil {
					return ""
				}
				return token
			},
			secret:   "secret",
			wantRole: "admin",
			wantErr:  false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			role, err := jwttoken.Check(c.token(), c.secret)
			assert.Equal(t, c.wantRole, role)
			if c.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
