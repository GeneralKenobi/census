package config

var defaultConfig = Config{
	Global: Global{
		ShutdownTimeoutSeconds: 30,
	},
	HttpServer: HttpServer{
		Port:                   8080,
		ShutdownTimeoutSeconds: 30,
	},
	Postgres: Postgres{
		Port:                  5432,
		DefaultTimeoutSeconds: 30,
	},
}