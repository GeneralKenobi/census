package config

var defaultConfig = Config{
	Global: Global{
		Database:               "postgres",
		ShutdownTimeoutSeconds: 30,
	},
	HttpServer: HttpServer{
		Port:                   8080,
		ShutdownTimeoutSeconds: 30,
	},
	Postgres: Postgres{
		Port:      5432,
		VerifyTls: true,
	},
	Mongo: Mongo{
		Port:           27017,
		TimeoutSeconds: 30,
		VerifyTls:      true,
	},
}
