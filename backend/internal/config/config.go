package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DBDriver     string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	DBURL        string
	RedisHost    string
	RedisPort    string
	RedisPass    string
	JWTSecret    string
	JWTExpiresIn int
	Environment  string
	CORSOrigin   string
}

func Load() *Config {
	_ = godotenv.Load() // Carrega .env se existir

	port := getEnv("PORT", "8080")
	dbDriver := getEnv("DB_DRIVER", "postgres") // "postgres" ou "sqlite"
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "sistema_agendamento")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")
	dbURL := getEnv("DATABASE_URL", "")

	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPass := getEnv("REDIS_PASSWORD", "")

	jwtSecret := getEnv("JWT_SECRET", "super-secret-scheduling-jwt-key-2026-secure-token")
	jwtExpiresStr := getEnv("JWT_EXPIRES_IN_HOURS", "72")
	jwtExpires, err := strconv.Atoi(jwtExpiresStr)
	if err != nil {
		jwtExpires = 72
	}

	env := getEnv("APP_ENV", "development")
	corsOrigin := getEnv("CORS_ORIGIN", "*")

	return &Config{
		Port:         port,
		DBDriver:     dbDriver,
		DBHost:       dbHost,
		DBPort:       dbPort,
		DBUser:       dbUser,
		DBPassword:   dbPassword,
		DBName:       dbName,
		DBSSLMode:    dbSSLMode,
		DBURL:        dbURL,
		RedisHost:    redisHost,
		RedisPort:    redisPort,
		RedisPass:    redisPass,
		JWTSecret:    jwtSecret,
		JWTExpiresIn: jwtExpires,
		Environment:  env,
		CORSOrigin:   corsOrigin,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
