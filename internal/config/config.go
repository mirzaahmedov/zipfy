package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Postgres struct {
		Host string `env:"POSTGRES_HOST" env-default:"localhost"`
		Port string `env:"POSTGRES_PORT" env-default:"5432"`
		User string `env:"POSTGRES_USER" env-default:"postgres"`
		PWD  string `env:"POSTGRES_PWD"`
		DB   string `env:"POSTGRES_DBNAME" env-default:"zipfy"`
		SSL  string `env:"POSTGRES_SSL" env-default:"disable"`
	}
	Redis struct {
		Addr string `env:"REDIS_ADDR" env-default:"localhost:6379"`
		PWD  string `env:"REDIS_PWD"`
		DB   int    `env:"REDIS_DB" env-default:"0"`
	}
	HTTP struct {
		Addr string `env:"HTTP_ADDR" env-default:":8080"`
	}
	JWT struct {
		Secret string `env:"JWT_SECRET" env-required:"true"`
	}
}

func Load() (*Config, error) {
	c := new(Config)

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file %v", err)
	}

	err = cleanenv.ReadEnv(c)
	if err != nil {
		return nil, fmt.Errorf("Error loading env variables %v", err)
	}

	return c, nil
}
