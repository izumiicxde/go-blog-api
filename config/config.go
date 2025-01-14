package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string `env:"PUBLIC_HOST"`
	Port       string `env:"PORT"`

	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBAddress  string `env:"DATABASE_URL"`
	DBName     string `env:"DB_NAME"`

	JWTSecret     string `env:"JWT_SECRET"`
	JWTExpiration int64  `env:"JWT_EXPIRATION"`
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	envs, err := env.ParseAs[Config]()
	if err != nil {
		log.Fatal(fmt.Errorf("error loading the env variables: %w", err))
		return Config{}
	}
	return envs
}
