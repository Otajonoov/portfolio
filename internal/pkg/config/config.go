package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Environment                string
	LogLevel                   string
	HttpPort                   string
	CtxTimeout                 int64
	Postgres                   Postgres
}

type Postgres struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresDatabase string
	PostgresPassword string
	DatabaseURL      string
}

func Load(path string) Config {
	godotenv.Load(path + "/.env")

	conf := viper.New()
	conf.AutomaticEnv()

	return Config{
		Environment:                conf.GetString("ENVIRONMENT"),
		CtxTimeout:                 conf.GetInt64("CTX_TIMEOUT"),
		LogLevel:                   conf.GetString("LOG_LEVEL"),
		HttpPort:                   conf.GetString("HTTP_PORT"),
		
		Postgres: Postgres{
			PostgresDatabase: conf.GetString("POSTGRES_DATABASE"),
			PostgresUser:     conf.GetString("POSTGRES_USER"),
			PostgresPassword: conf.GetString("POSTGRES_PASSWORD"),
			PostgresPort:     conf.GetString("POSTGRES_PORT"),
			PostgresHost:     conf.GetString("POSTGRES_HOST"),
			DatabaseURL:      conf.GetString("DATABASE_URL"),
		},
	}
}
