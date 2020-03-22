package config

//TestConfig - вспомогательная функция для получения конфигурации
func TestConfig() *Config {
	return &Config{
		Server: &Server{
			Port:   33333,
			Secret: "k8Lds17DDl100,ask",
		},
		Db: &Db{
			Driver:           "postgres",
			ConnectionString: "postgres://postgres:postgres@localhost:5432/test_authorization?sslmode=disable",
		},
	}
}