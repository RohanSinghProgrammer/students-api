package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Address string `yaml:"address" default:"4000"`
}

type Config struct {
	Env		  string `yaml:"env"  env-require:"true"`
	StoragePath string `yaml:"storage_path" env-require:"true"`
	HTTPServer `yaml:"http_server"`
}

func MustLoadConfig() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "Path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("CONFIG_PATH environment variable or --config flag must be set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Configuration file does not exist at path: %s", configPath)
	}

	var config Config

	err := cleanenv.ReadConfig(configPath, &config)

	if err != nil {
		log.Fatalf("Failed to read configuration file: %s", err.Error())
	}

	return &config
}