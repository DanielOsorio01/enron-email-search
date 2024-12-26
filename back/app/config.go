package app

type Config struct {
	dbAddr     string
	dbUser     string
	dbPassword string
	serverPort uint16
}

func LoadConfig() *Config {
	return &Config{
		dbAddr:     "http://localhost:4080",
		dbUser:     "admin",
		dbPassword: "Complexpass#123",
		serverPort: 3000,
	}

	// TODO: Load configuration from environment variables
}
