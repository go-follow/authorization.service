package store

import "github.com/go-follow/authorization.service/config"

//TestConfigDB - helper для получения конфигурации DB
func TestConfigDB() *config.Db {
	return &config.Db{
		Driver:           "postgres",
		ConnectionString: "postgres://postgres:postgres@localhost:5432/test_authorization?sslmode=disable",
	}
}