package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port                    string
	LinuxOdPriceUrl         string
	LinuxOdPricePreviousUrl string
	GinEnv                  string
}

func Load(envFile string) (*Config, error) {
	if envFile == "" {
		envFile = ".env"
	}

	fmt.Printf("Loading %s file\n", envFile)

	godotenv.Load(envFile)

	linuxOdPriceUrl := os.Getenv("LINUX_ON_DEMAND_PRICE_URL")
	if linuxOdPriceUrl == "" {
		return nil, fmt.Errorf("Missing LINUX_ON_DEMAND_PRICE_URL")
	}

	linuxOdPricePreviousUrl := os.Getenv("LINUX_ON_DEMAND_PRICE_PREVIOUS_GEN_URL")
	if linuxOdPricePreviousUrl == "" {
		return nil, fmt.Errorf("Missing LINUX_ON_DEMAND_PRICE_PREVIOUS_GEN_URL")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	ginEnv := os.Getenv("GIN_ENV")
	if ginEnv == "" {
		ginEnv = "development"
	}

	return &Config{
		port,
		linuxOdPriceUrl,
		linuxOdPricePreviousUrl,
		ginEnv,
	}, nil
}
