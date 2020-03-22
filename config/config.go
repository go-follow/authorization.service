package config

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
