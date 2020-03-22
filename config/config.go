package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

//Config - конфигурационные настройки для всего проекта
type Config struct {
	Server *Server
	Db     *Db
}

//Server - настройки для запуска рест сервера
type Server struct {
	Port   uint16
	Secret string	
}

//Db - настройки к подключению Db
type Db struct {
	Driver           string
	ConnectionString string
}

//ReadConfig - чтение всех настроек из enviroment
func ReadConfig() (*Config, error) {

	server, err := serverConfig()
	if err != nil {
		return nil, err
	}
	db, err := dbConfig()
	if err != nil {
		return nil, err
	}
	return &Config{
		Server: server,
		Db:     db,
	}, nil
}

func dbConfig() (*Db, error) {
	driver := strings.TrimSpace(os.Getenv("DB_DRIVER"))
	connectionString := strings.TrimSpace(os.Getenv("DB_CONNECTION_STRING"))
	if driver == "" {
		return nil, fmt.Errorf("empty env DB_DRIVER")
	}
	if connectionString == "" {
		return nil, fmt.Errorf("empty env DB_CONNECTION_STRING")
	}
	return &Db{
		Driver:           driver,
		ConnectionString: connectionString,
	}, nil
}

func serverConfig() (*Server, error) {
	portS := strings.TrimSpace(os.Getenv("SERVER_PORT"))
	if portS == "" {
		return nil, fmt.Errorf("empty env SERVER_PORT")
	}

	p, err := strconv.ParseUint(portS, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("env SERVER_PORT should be integer and have range from 0 to  65535: %v", err)
	}

	secret := strings.TrimSpace(os.Getenv("SERVER_SECRET"))
	if secret == "" {
		return nil, fmt.Errorf("empty env SERVER_SECRET")
	}
	return &Server{
		Port:   uint16(p),
		Secret: secret,
	}, nil
}
