package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Database  DatabaseConfig
	Redis     RedisConfig
	Server    ServerConfig
	JWT       JWTConfig
	Auth      AuthConfig
	Cloudinary CloudinaryConfig
	Env       string
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
	Secret             string
	AccessTokenMinutes int
}

type AuthConfig struct {
	RefreshTokenDays      int
	RefreshCookieName     string
	RefreshCookiePath     string
	RefreshCookieDomain   string
	RefreshCookieSecure   bool
	RefreshCookieSameSite string
}

type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
}

func Load() *Config {
	// Загружаем .env файл если он существует
	godotenv.Load()

	// Парсим DATABASE_URL если он есть, иначе используем отдельные переменные
	dbConfig := parseDatabaseConfig()

	return &Config{
		Database: dbConfig,
		Redis: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     getEnvAsInt("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		Server: ServerConfig{
			Host: getEnvWithDefault("SERVER_HOST", "0.0.0.0"),
			Port: getServerPort(),
		},
		JWT: JWTConfig{
			Secret:             getEnvWithDefault("JWT_SECRET", "your-secret-key-change-in-production"),
			AccessTokenMinutes: getEnvAsIntWithDefault("ACCESS_TOKEN_MINUTES", 15),
		},
		Auth: AuthConfig{
			RefreshTokenDays:      getEnvAsIntWithDefault("REFRESH_TOKEN_DAYS", 30),
			RefreshCookieName:     getEnvWithDefault("REFRESH_COOKIE_NAME", "refresh_token"),
			RefreshCookiePath:     getEnvWithDefault("REFRESH_COOKIE_PATH", "/"),
			RefreshCookieDomain:   os.Getenv("COOKIE_DOMAIN"),
			RefreshCookieSecure:   getEnvWithDefault("ENV", "development") == "production",
			RefreshCookieSameSite: getEnvWithDefault("COOKIE_SAMESITE", "Lax"),
		},
		Cloudinary: CloudinaryConfig{
			CloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
			APIKey:    os.Getenv("CLOUDINARY_API_KEY"),
			APISecret: os.Getenv("CLOUDINARY_API_SECRET"),
		},
		Env: getEnvWithDefault("ENV", "development"),
	}
}

func parseDatabaseConfig() DatabaseConfig {
	// Сначала проверяем DATABASE_URL
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		return parseDatabaseURL(databaseURL)
	}

	// Если DATABASE_URL нет, используем отдельные переменные
	return DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     getEnvAsInt("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}

func parseDatabaseURL(databaseURL string) DatabaseConfig {
	// Убираем префикс postgresql://
	databaseURL = strings.TrimPrefix(databaseURL, "postgresql://")

	// Парсим URL: user:password@host:port/database?sslmode=require
	parts := strings.Split(databaseURL, "@")
	if len(parts) != 2 {
		// Если формат неправильный, возвращаем пустую конфигурацию
		return DatabaseConfig{}
	}

	// Парсим user:password
	userPass := strings.Split(parts[0], ":")
	user := userPass[0]
	password := ""
	if len(userPass) > 1 {
		password = userPass[1]
	}

	// Парсим host:port/database?sslmode=require
	hostDB := strings.Split(parts[1], "?")[0] // убираем параметры
	hostPortDB := strings.Split(hostDB, "/")
	if len(hostPortDB) != 2 {
		return DatabaseConfig{}
	}

	database := hostPortDB[1]
	hostPort := strings.Split(hostPortDB[0], ":")
	host := hostPort[0]
	port := 5432 // дефолтный порт PostgreSQL
	if len(hostPort) > 1 {
		if p, err := strconv.Atoi(hostPort[1]); err == nil {
			port = p
		}
	}

	return DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     database,
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

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getServerPort() int {
	// Render автоматически устанавливает переменную PORT
	if port := os.Getenv("PORT"); port != "" {
		if intValue, err := strconv.Atoi(port); err == nil {
			return intValue
		}
	}
	// Если PORT не установлен, используем SERVER_PORT
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if intValue, err := strconv.Atoi(port); err == nil {
			return intValue
		}
	}
	// Дефолтный порт
	return 8080
}
