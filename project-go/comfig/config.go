package config

import (
    "os"
)

type Config struct {
    DBPath    string
    JWTSecret string
    Port      string
}

func LoadConfig() (*Config, error) {
    return &Config{
        DBPath:    getEnv("DB_PATH", "./urlshortener.db"),
        JWTSecret: getEnv("JWT_SECRET", "default-secret-key-change-me"),
        Port:      getEnv("PORT", "3000"),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
