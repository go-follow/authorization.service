package db

import (
	"fmt"

	"github.com/go-follow/authorization.service/config"
	"github.com/jmoiron/sqlx"

	//Официальный Драйвер postgres
	_ "github.com/lib/pq"
)

//Connect - открытие соединения с БД
func Connect(d *config.Db) (*sqlx.DB, error) {
	if d == nil {
		return nil, fmt.Errorf("empty config for db")
	}
	c, err := sqlx.Open(d.Driver, d.ConnectionString)
	if err != nil {
		return nil, err
	}
	if err := c.Ping(); err != nil {
		return nil, err
	}
	return c, nil
}
