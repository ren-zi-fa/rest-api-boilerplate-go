package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost                    string
	Port                          string
	DBUser                        string
	DBPassword                    string
	DBAddress                     string
	DBName                        string
	JWTSecret                     string
	REFRESH_TOKEN_EXPIRE_DURATION time.Duration
	ACCESS_TOKEN_EXPIRE_DURATION  time.Duration
}

var Envs = initConfig()

func initConfig() *Config {
	env := os.Getenv("APP_ENV") 

	if env == "prod" {
		godotenv.Load(".env.prod")
	} else {
		godotenv.Load(".env.dev")
	}

	return &Config{
		PublicHost:                    getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                          getEnv("PORT", "8080"),
		DBUser:                        getEnv("DB_USER", "root"),
		DBPassword:                    getEnv("DB_PASSWORD", ""),
		DBAddress:                     fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:                        getEnv("DB_NAME", "my_db"),
		JWTSecret:                     getEnv("JWT_SECRET", "your-256-bit-secret"),
		REFRESH_TOKEN_EXPIRE_DURATION: getEnvDuration("REFRESH_TOKEN_EXPIRE_DURATION", time.Hour*24*7),
		ACCESS_TOKEN_EXPIRE_DURATION:  getEnvDuration("ACCESS_TOKEN_EXPIRE_DURATION", 15*time.Minute),
	}

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback

}
func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		d, err := time.ParseDuration(value)
		if err != nil {
			return fallback
		}
		return d
	}

	return fallback

}
