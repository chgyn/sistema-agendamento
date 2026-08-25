package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DBDriver           string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	DBURL              string
	RedisHost          string
	RedisPort          string
	RedisPass          string
	JWTSecret          string
	JWTExpiresIn       int
	Environment        string
	CORSOrigin         string
	AsaasAPIKey        string
	AsaasBaseURL       string
	AsaasWebhookSecret string
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

	asaasAPIKey := getEnv("ASAAS_API_KEY", "$aact_hmlg_000MzkwODA2MWY2OGM3MWRlMDU2NWM3MzJlNzZmNGZhZGY6OjQ1NjY2ZDQ4LWQxM2YtNDA1YS1hMmRjLTE0MjQ5NTQwNzJhNTo6JGFhY2hfOTVkODc0YmYtOWIyYi00YWEwLWFlZjYtMjgxY2YwYmQ1OTQy")
	asaasBaseURL := getEnv("ASAAS_BASE_URL", "https://sandbox.asaas.com/api/v3")
	asaasWebhookSecret := getEnv("ASAAS_WEBHOOK_SECRET", "whsec_hhwwxdTmKX30aXHwqiSGB9jtv9zYbM8IAo-oCDNeSdg")

	return &Config{
		Port:               port,
		DBDriver:           dbDriver,
		DBHost:             dbHost,
		DBPort:             dbPort,
		DBUser:             dbUser,
		DBPassword:         dbPassword,
		DBName:             dbName,
		DBSSLMode:          dbSSLMode,
		DBURL:              dbURL,
		RedisHost:          redisHost,
		RedisPort:          redisPort,
		RedisPass:          redisPass,
		JWTSecret:          jwtSecret,
		JWTExpiresIn:       jwtExpires,
		Environment:        env,
		CORSOrigin:         corsOrigin,
		AsaasAPIKey:        asaasAPIKey,
		AsaasBaseURL:       asaasBaseURL,
		AsaasWebhookSecret: asaasWebhookSecret,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
