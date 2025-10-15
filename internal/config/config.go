package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	Server   ServerConfig
	JWT      JWTConfig
	Env      string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

type ServerConfig struct {
	Host string
	Port int
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

func Load() *Config {
	// Загружаем .env файл если он существует
	godotenv.Load()

	return &Config{
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     getEnvAsInt("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     getEnvAsInt("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		Server: ServerConfig{
			Host: os.Getenv("SERVER_HOST"),
			Port: getEnvAsInt("SERVER_PORT"),
		},
		JWT: JWTConfig{
			Secret:      os.Getenv("JWT_SECRET"),
			ExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS"),
		},
		Env: os.Getenv("ENV"),
	}
}

func getEnvAsInt(key string) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return 0
}