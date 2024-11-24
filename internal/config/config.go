package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/reversersed/zyfra/pkg/mongo"
)

type Config struct {
	Server   *ServerConfig
	Database *mongo.DatabaseConfig
}

type ServerConfig struct {
	Host string `env:"SERVER_HOST" env-required:"true" env-description:"Server listening address"`
	Port int    `env:"SERVER_PORT" env-required:"true" env-description:"Server listening port"`
}

func GetConfig(path string) (*Config, error) {
	server := &ServerConfig{}
	database := &mongo.DatabaseConfig{}

	if err := cleanenv.ReadConfig(path, server); err != nil {
		var header string = "Server part config"
		desc, _ := cleanenv.GetDescription(server, &header)
		return nil, fmt.Errorf("%v: %s", err, desc)
	}
	if err := cleanenv.ReadConfig(path, database); err != nil {
		var header string = "Database part config"
		desc, _ := cleanenv.GetDescription(database, &header)
		return nil, fmt.Errorf("%v\n%s", err, desc)
	}
	config := &Config{
		Server:   server,
		Database: database,
	}

	return config, nil
}
