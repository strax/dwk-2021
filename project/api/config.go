package main

import (
	"fmt"
	"os"
	"strconv"
)

func MustLookupEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Environment value %s not is not present", value))
	}
	return value
}

type AppConfig struct {
	PathPrefix string
	NatsURL    string
	DBConfig
}

type DBConfig struct {
	Host     string
	Port     uint16
	User     string
	Password string
	Database string
}

func AppConfigFromEnv() AppConfig {
	dbConfig := DBConfigFromEnv()
	pathPrefix, _ := os.LookupEnv("PATH_PREFIX")
	natsURL := MustLookupEnv("NATS_URL")
	return AppConfig{pathPrefix, natsURL, dbConfig}
}

func DBConfigFromEnv() DBConfig {
	host := MustLookupEnv("PGHOST")
	password := MustLookupEnv("PGPASS")
	port, err := strconv.ParseInt(MustLookupEnv("PGPORT"), 10, 16)
	if err != nil {
		panic("Invalid PGPORT")
	}
	user := MustLookupEnv("PGUSER")
	database := MustLookupEnv("PGDATABASE")
	return DBConfig{
		Host:     host,
		Password: password,
		Port:     uint16(port),
		User:     user,
		Database: database,
	}
}

func (c *DBConfig) ToPostgresConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", c.User, c.Password, c.Host, c.Port, c.Database)
}
