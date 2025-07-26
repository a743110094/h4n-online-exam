package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port           string
	DatabaseURL    string
	RedisURL       string
	JWTSecret      string
	AIAPIKey       string
	AIURL          string
}

var config *Config
var AppConfig *Config

func Init() {
	config = &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    buildDatabaseURL(),
		RedisURL:       getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:      getEnv("JWT_SECRET", "online-exam-system-jwt-secret-key-2024"),
		AIAPIKey:       getEnv("AI_API_KEY", ""),
		AIURL:          getEnv("AI_URL", "https://api.openai.com/v1/chat/completions"),
	}
	AppConfig = config
}

func buildDatabaseURL() string {
	// 如果设置了DATABASE_URL，直接使用
	if dbURL := getEnv("DATABASE_URL", ""); dbURL != "" {
		return dbURL
	}
	
	// 否则从分别的环境变量构建
	dbName := getEnv("DB_NAME", "")
	
	// 如果DB_NAME以.db结尾，使用SQLite
	if dbName != "" && (dbName[len(dbName)-3:] == ".db" || dbName == ":memory:") {
		return dbName
	}
	
	// 构建PostgreSQL连接字符串
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	if dbName == "" {
		dbName = "online_exam_system"
	}
	sslmode := getEnv("DB_SSLMODE", "disable")
	
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbName, sslmode)
}

func GetConfig() *Config {
	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}