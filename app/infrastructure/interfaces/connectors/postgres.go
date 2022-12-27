package connectors

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// postgresDBInstance is an instance for of postgres db in gorm format
	postgresDBInstance *gorm.DB
	postgresHost       string
	postgresPort       string
	postgresUser       string
	postgresPassword   string
	postgresDB         string
	postgresDBLogMode  bool
)

// InitPostgresDB is a init function for setting config for postgres
func InitPostgresDB(config DatabaseConfig) {
	postgresHost = config.Host
	postgresPort = config.Port
	postgresUser = config.User
	postgresPassword = config.Password
	postgresDB = config.DB
	postgresDBLogMode = config.DBLogMode
}

// ConnectPostgresDB is a connector function for connecting postgres
func ConnectPostgresDB() *gorm.DB {
	return connect(postgresDBLogMode)
}

// PostgresConnection is a postgres connection string
func PostgresConnection() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		postgresHost,
		postgresPort,
		postgresUser,
		postgresPassword,
		postgresDB,
	)
}

func connect(dbLogMode bool) *gorm.DB {
	var err error
	connection := PostgresConnection()
	mut.Lock()
	defer mut.Unlock()
	var logLevel logger.LogLevel
	if dbLogMode {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}
	if postgresDBInstance == nil {
		postgresDBInstance, err = gorm.Open(postgres.Open(connection), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})
		if err != nil {
			panic(err.Error())
		}
	}
	return postgresDBInstance
}
