package db_test

import (
	"testing"

	"github.com/go-follow/authorization.service/config"
	"github.com/go-follow/authorization.service/internal/service/store"
	"github.com/go-follow/authorization.service/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	cases := []struct {
		name    string
		cfg     func() *config.Db
		wantErr bool
	}{
		{
			name:    "empty_config",
			cfg:     func() *config.Db { return nil },
			wantErr: true,
		},
		{
			name: "error_open",
			cfg: func() *config.Db {
				db := store.TestConfigDB()
				db.Driver = "unknown"
				return db
			},
			wantErr: true,
		},
		{
			name: "error_ping",
			cfg: func() *config.Db {
				db := store.TestConfigDB()
				db.ConnectionString = ""
				return db
			},
			wantErr: true,
		},
		{
			name: "success_connect",
			cfg: func() *config.Db {
				return store.TestConfigDB()
			},
			wantErr: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			conn, err := db.Connect(c.cfg())
			if c.wantErr {
				assert.Error(t, err)
				assert.Nil(t, conn)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, conn)
			}
		})
	}
}