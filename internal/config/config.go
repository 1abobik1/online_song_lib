package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env            string `env:"ENV" env-default:"development"`
	StoragePath    string `env:"DATABASE_URL" env-required:"true"`
	HTTPServer     string `env:"HTTP_SERVER_ADDRESS"`
	ExternalApiURL string `env:"EXTERNAL_API_URL"`
}

func MustLoad() *Config {
	path := getConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist:" + path)
	}

	// Загружаем переменные из указанного файла
	err := godotenv.Load(path)
	if err != nil {
		panic(fmt.Sprintf("No .env file found at %s, relying on environment variables", path))
	}

	var cfg Config

	// Загружаем переменные окружения в структуру
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(fmt.Sprintf("Failed to load environment variables: %v", err))
	}

	return &cfg
}

func getConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	return res
}
