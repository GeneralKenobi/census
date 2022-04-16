package config

// Config aggregates configuration for modules.
type Config struct {
	Global     Global     `json:"global"`
	HttpServer HttpServer `json:"httpServer"`
	Postgres   Postgres   `json:"postgres"`
	Mongo      Mongo      `json:"mongo"`
}

// Global contains general configuration or configuration for the entire application.
type Global struct {
	Database               string `json:"database"`               // Which database to use, currently supported are postgres and mongo
	ShutdownTimeoutSeconds int    `json:"shutdownTimeoutSeconds"` // Maximum time for graceful shutdown of the application
}

type HttpServer struct {
	Port                   int    `json:"port"`                   // Port to listen on
	ShutdownTimeoutSeconds int    `json:"shutdownTimeoutSeconds"` // Graceful shutdown time
	Tls                    bool   `json:"tls"`                    // Flag to enable/disable TLS
	TlsCertPath            string `json:"tlsCertPath"`            // Path to the TLS certificate to serve
	TlsCertKeyPath         string `json:"tlsCertKeyPath"`         // Path to the key corresponding to TlsCertPath
}

type Postgres struct {
	Host      string `json:"host"`      // DB server host, e.g. my-postgres.com or 10.101.146.170
	Port      int    `json:"port"`      // Port the DB is listening on
	User      string `json:"user"`      // User to use for connecting to the DB
	Password  string `json:"password"`  // Password authenticating User
	Database  string `json:"database"`  // Database to use
	VerifyTls bool   `json:"verifyTls"` // Flag to turn TLS/SSL verification on/off
}

type Mongo struct {
	Host           string `json:"host"`           // DB server host, e.g. my-mongo.com or 10.101.146.170
	Port           int    `json:"port"`           // Port the DB is listening on
	User           string `json:"user"`           // User to use for connecting to the DB
	Password       string `json:"password"`       // Password authenticating User
	ReplicaSet     string `json:"replicaSet"`     // Name of the replica set to connect to
	Database       string `json:"database"`       // Database to use
	TimeoutSeconds int    `json:"timeoutSeconds"` // Timeout to use for opening/closing DB connection
	VerifyTls      bool   `json:"verifyTls"`      // Flag to turn TLS/SSL verification on/off
}
