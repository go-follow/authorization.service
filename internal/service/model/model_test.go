package model_test

import (
	"testing"

	"github.com/go-follow/authorization.service/internal/service/model"
	"github.com/stretchr/testify/assert"
)

func TestUserValidate(t *testing.T) {
	cases := []struct {
		name    string
		u       func() *model.User
		wantErr bool
	}{
		{
			name: "empty_name",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Name = "  "
				return u
			},
			wantErr: true,
		},
		{
			name: "length_name",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Name = model.TestGenerateString(t, 31)
				return u
			},
			wantErr: true,
		},
		{
			name: "empty_password",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = ""
				return u
			},
			wantErr: true,
		},
		{
			name: "length_password",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = model.TestGenerateString(t, 5)
				return u
			},
			wantErr: true,
		},
		{
			name: "empty_password_hash",
			u: func() *model.User {
				u := model.TestUser(t)
				u.PasswordEncode = "   "
				return u
			},
			wantErr: true,
		},
		{
			name: "length_password_hash",
			u: func() *model.User {
				u := model.TestUser(t)
				u.PasswordEncode = model.TestGenerateString(t, 256)
				return u
			},
			wantErr: true,
		},
		{
			name: "empty_role",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Role = "   "
				return u
			},
			wantErr: true,
		},
		{
			name: "not_valid_role",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Role = "unknown"
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
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.u().Validate()
			if c.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSetHashPassword(t *testing.T) {
	//password := "k8sAd9jd"
	password := model.TestGenerateString(t, 999)
	secret := "falkKKSlk2d3:kAl'83"
	cases := []struct {
		name    string
		u       *model.User
		wantErr bool
	}{
		{
			name: "success",
			u: &model.User{
				Password:       password,
				PasswordEncode: "",
			},
			wantErr: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.u.SetHashPassword(secret)
			if c.wantErr {
				assert.Error(t, err)
				assert.Empty(t, c.u.PasswordEncode)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, c.u.PasswordEncode)
			}
		})
	}
}

func TestComparePassword(t *testing.T) {
	password := "k8sAd9jc"
	secret := "falkKKSlk2d3:kAl'83"
	//hashPassword := "e0c02ae934ef2864544a0e5144cb03d9c4e9263ff83a0760b7054977792cf7fee8ccad630eda50b4cab052a9fc03bb2d1bc2f0764a99c9555299f228d3f32bff"
	cases := []struct {
		name          string
		u             func() *model.User
		secret        string
		compareResult bool
	}{
		{
			name: "empty_password",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = "   "
				if err := u.SetHashPassword(secret); err != nil {
					return &model.User{}
				}
				return u
			},
			secret:        secret,
			compareResult: false,
		},
		{
			name: "empty_secret",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = password
				if err := u.SetHashPassword("   "); err != nil {
					return &model.User{}
				}
				return u
			},
			secret:        secret,
			compareResult: false,
		},
		{
			name: "success",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = password
				if err := u.SetHashPassword(secret); err != nil {
					return &model.User{}
				}
				return u
			},
			secret:        secret,
			compareResult: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.u().ComparePassword(c.secret)
			assert.Equal(t, c.compareResult, result)
		})
	}
}
