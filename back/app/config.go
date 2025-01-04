package app

import (
	"os"
	"strconv"
)

type Config struct {
	dbAddr     string
	dbUser     string
	dbPassword string
	serverPort uint16
}

func LoadConfig() *Config {
	port, err := strconv.ParseUint(os.Getenv("SERVER_PORT"), 10, 16)
	if err != nil {
		port = 3000
	}

	return &Config{
		dbAddr:     getEnv("DB_HOST", "http://localhost:4080"),
		dbUser:     getEnv("DB_USER", "admin"),
		dbPassword: getEnv("DB_PASSWORD", "Complexpass#123"),
		serverPort: uint16(port),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
