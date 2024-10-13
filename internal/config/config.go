package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server" `
}

// It will serialize the cofig it will get
func MustLoad() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	//Suppose if config is passed through config flag
	if configPath == "" {
		flags := flag.String("config", "", "path to the config file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("config path is not set")
		}
	}

	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	err = cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("can not read config file: %s", err.Error())
	}

	return &cfg
}
