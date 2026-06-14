package config

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	ProjectRoot string
	BaseUrl     string
	LogLevel   string
    LogFileName string
    LogDir    string
    ServiceName string

}

func LoadConfig(serviceName, env string) (*Config, error) {
	envFile := fmt.Sprintf("../../envs/.env.%s", env) // change 
	

	if err := godotenv.Load(envFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("ENV file isn't present, continuing without it:", envFile)
		} else {
			return nil, fmt.Errorf("failed to load %s: %w", envFile, err)
		}
	}

	config := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		ProjectRoot: getEnv("PROJECT_ROOT", "" ), // change 
		BaseUrl:     getEnv("BASE_URL", ""),
		LogLevel: getEnv("LOG_LEVEL",""),
        LogFileName: getEnv("LOG_FILE_NAME",""),
        LogDir: getEnv("LOG_DIR",""),
        ServiceName: serviceName,
	}

	if config.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

